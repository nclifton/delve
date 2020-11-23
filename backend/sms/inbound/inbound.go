package inbound

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/middleware/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/recovery"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type InboundOptions struct {
	NrApp      func(http.Handler) http.Handler
	SMSClient  *rpc.Client
	BuuRestUrl string
}

// API wraps an instance of our api app
type InboundAPI struct {
	opts   *InboundOptions
	router *httprouter.Router
	sms    *rpc.Client
}

// Handler exposes the router
func (a *InboundAPI) Handler() http.Handler {
	return a.router
}

// New creates our api "app", i.e. the http handler
func NewInboundAPI(opts *InboundOptions) *InboundAPI {
	api := &InboundAPI{
		opts: opts,
		sms:  opts.SMSClient,
	}

	loggerM := logger.New(&logger.Options{
		Verbose:    true,
		UseXRealIp: true,
	})

	newrelicM := opts.NrApp

	// maybe not needed with httprouter panic handler
	recoveryM := recovery.New(&recovery.Options{
		PrintStack: true,
	})

	// define the middleware chains for our api endpoints
	// we can group them and expand chains
	baseChain := alice.New(loggerM, newrelicM, recoveryM)

	// configure routes
	router := httprouter.New() // /v1/
	api.router = router

	// we also need a route and chain for 404
	router.NotFound = NewPlainRoute(api, baseChain, NotFoundRoute)
	router.POST("/v1/sms/dlr", NewRoute(api, baseChain, InboundDLRPOST))
	//router.GET("/v1/sms/mo", NewRoute(api, baseChain, InboundMOGET))

	return api
}

// Route is some bullshit type so we can make httprouter chain with alice
// and still access params and such
type Route struct {
	params   httprouter.Params
	w        http.ResponseWriter
	r        *http.Request
	endpoint func(*Route)
	api      *InboundAPI
}

// NotFoundRoute route for convenience
func NotFoundRoute(r *Route) {
	http.Error(r.w, "Route not found", http.StatusNotFound)
}

// NewRoute returns a handler wrapped in middleware
// it's super important that route objects are created when requests come in
// if they are created before, the same object is shared for all requests
// and you fuck yourself with concurrent access... derp!
func NewRoute(api *InboundAPI, chain alice.Chain, endpoint func(*Route)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		f := &Route{
			params:   p,
			endpoint: endpoint,
			api:      api,
		}

		chain.Then(f).ServeHTTP(w, r)
	}
}

// NewPlainRoute for plain http.Handler's
func NewPlainRoute(api *InboundAPI, chain alice.Chain, endpoint func(*Route)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := &Route{
			endpoint: endpoint,
			api:      api,
		}

		chain.Then(f).ServeHTTP(w, r)
	}
}

// ServeHTTP for accepting wrapped objects
func (r *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.w = w
	r.r = req
	r.endpoint(r)
}

func InboundDLRPOST(r *Route) {
	log.Print("Got DLR")
	err := r.r.ParseForm()
	if err != nil {
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to parse params: %s", err)
		return
	}

	log.Printf("DLR Form values: %+v", r.r.Form)

	if len(r.r.Form["msgid"]) < 1 || len(r.r.Form["state"]) < 1 {
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to parse DLR response from Alaris")
		return
	}

	log.Printf("[%s] DLR State: %s", r.r.FormValue("msgid"), r.r.FormValue("state"))

	time, err := time.Parse(time.RFC3339, r.r.FormValue("time"))
	if err != nil {
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to parse DLR response from Alaris: %s", err)
	}

	err = r.api.sms.QueueDLR(rpc.QueueDLRParams{
		MessageID:  r.r.FormValue("msgid"),
		State:      r.r.FormValue("state"),
		To:         r.r.FormValue("to"),
		Time:       time,
		ReasonCode: r.r.FormValue("reasoncode"),
		MCC:        r.r.FormValue("mcc"),
		MNC:        r.r.FormValue("mnc"),
	})
	if err != nil {
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to queue DLR response from Alaris: %s", err)
	}

	r.w.WriteHeader(http.StatusOK)
	fmt.Fprintf(r.w, "0 OK")
}

/*func InboundMOGET(r *Route) {
	err := r.r.ParseForm()
	if err != nil {
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to parse params: %s", err)
		return
	}

	if len(r.r.Form["source"]) < 1 || len(r.r.Form["message"]) < 1 || len(r.r.Form["dest"]) < 1 {
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to parse MO request from Kannel")
		return
	}

	source := r.r.Form["source"][0]
	dest := r.r.Form["dest"][0]
	message := r.r.Form["message"][0]

	log.Printf("MO: source:%s dest:%s message:%s", source, dest, message)

	// Is this a response to an sms we have sent
	sms, err := r.api.sms.FindByMO(source, dest)
	if err != nil {
		// if not redirect it to buu
		if err.Error() == mongo.ErrNoDocuments.Error() {
			log.Printf("Did not find related SMS so letting buu handle it: %s", err)
			redirectURL := fmt.Sprintf("%s/sms/mo?%s", r.api.opts.BuuRestUrl, r.r.URL.RawQuery)
			http.Redirect(r.w, r.r, redirectURL, http.StatusFound)
			return
		}
		log.Printf("Failed to find reply MO: %s", err)
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Found related sms: %+v", sms.SMS)

	isValid, err := r.api.sms.ValidateSender(dest, "", sms.SMS.AccountID)
	if err != nil {
		log.Printf("Failed to Validate Sender ID: %s", err)
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isValid.ValidAccount {
		log.Printf("Did not find related SMS so letting buu handle it: %s", err)
		redirectURL := fmt.Sprintf("%s/sms/mo?%s", r.api.opts.BuuRestUrl, r.r.URL.RawQuery)
		http.Redirect(r.w, r.r, redirectURL, http.StatusFound)
		return
	}

	// Does this account still own this sender id
	opts := rabbit.PublishOptions{
		RouteKey:     msg.MOProcess.RouteKey,
		Exchange:     msg.MOProcess.Exchange,
		ExchangeType: msg.MOProcess.ExchangeType,
	}

	mojob := msg.MOProcessSpec{
		ReplyTo:   sms.SMS.ID,
		AccountID: sms.SMS.AccountID,
		Source:    source,
		Dest:      dest,
		Message:   message,
	}

	err = r.api.db.RabbitPublish(opts, mojob)
	if err != nil {
		http.Error(r.w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to publish MO to queue response(%s) for message(%s) from Kannel: %s", message, sms.SMS.ID.Hex(), err)

		return
	}

	r.w.WriteHeader(http.StatusOK)
	fmt.Fprintf(r.w, "0 OK")
}*/
