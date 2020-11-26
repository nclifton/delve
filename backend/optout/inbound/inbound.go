package inbound

import (
	"log"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/middleware/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/recovery"
	rpc "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type InboundOptions struct {
	NrApp        func(http.Handler) http.Handler
	OptOutClient *rpc.Client
}

// API wraps an instance of our api app
type InboundAPI struct {
	opts   *InboundOptions
	router *httprouter.Router
	optout *rpc.Client
}

// Handler exposes the router
func (a *InboundAPI) Handler() http.Handler {
	return a.router
}

// New creates our api "app", i.e. the http handler
func NewInboundAPI(opts *InboundOptions) *InboundAPI {
	api := &InboundAPI{
		opts:   opts,
		optout: opts.OptOutClient,
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
	router.GET("/:linkID", NewRoute(api, baseChain, OptOutGET))
	router.POST("/:linkID", NewRoute(api, baseChain, OptOutPOST))

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

func OptOutGET(r *Route) {
	r.w.WriteHeader(http.StatusOK)

	optout, err := r.api.optout.FindByLinkID(rpc.FindByLinkIDParams{LinkID: r.params.ByName("linkID")})
	if err != nil {
		log.Printf("Err: %s", err.Error())
		return
	}
	err = renderLink(r.w, OptOutLink{
		Sender: optout.Sender,
		Link:   optout.LinkID,
	})
	if err != nil {
		log.Printf("Err: %s", err.Error())
	}
}

func OptOutPOST(r *Route) {
	r.w.WriteHeader(http.StatusOK)
	optout, err := r.api.optout.OptOutViaLink(rpc.OptOutViaLinkParams{LinkID: r.params.ByName("linkID")})
	if err != nil {
		log.Printf("Err: %s", err.Error())
		return
	}
	err = renderUnsubscribed(r.w, OptOutLink{
		Sender: optout.Sender,
	})
	if err != nil {
		log.Printf("Err: %s", err.Error())
	}
}
