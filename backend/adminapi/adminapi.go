package adminapi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/recovery"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
)

type AdminAPIOptions struct {
	NrApp         func(http.Handler) http.Handler
	AccountClient *account.Client
	SMSClient     *sms.Client
	MMSClient     *mms.Client
	SenderClient  senderpb.ServiceClient
}

type RPCClients struct {
	account *account.Client
	sms     *sms.Client
	mms     *mms.Client
	sender  senderpb.ServiceClient
}

// API wraps an instance of our api app
type AdminAPI struct {
	opts   *AdminAPIOptions
	router *httprouter.Router
	RPCClients
}

// Handler exposes the router
func (a *AdminAPI) Handler() http.Handler {
	return a.router
}

// New creates our api "app", i.e. the http handler
func NewAdminAPI(opts *AdminAPIOptions) *AdminAPI {
	clients := RPCClients{
		account: opts.AccountClient,
		sms:     opts.SMSClient,
		mms:     opts.MMSClient,
		sender:  opts.SenderClient,
	}

	api := &AdminAPI{
		opts:       opts,
		RPCClients: clients,
	}

	loggerM := logger.New(&logger.Options{
		Verbose:    true,
		UseXRealIp: true,
	})

	// maybe not needed with httprouter panic handler
	recoveryM := recovery.New(&recovery.Options{
		PrintStack: true,
	})

	// define the middleware chains for our api endpoints
	// we can group them and expand chains
	baseChain := alice.New(loggerM, recoveryM)

	// configure routes
	router := httprouter.New() // /v1/
	api.router = router

	if opts.NrApp != nil {
		baseChain.Append(opts.NrApp)
	}

	// we also need a route and chain for 404
	router.NotFound = NewPlainRoute(api, baseChain, NotFoundRoute)
	router.GET("/v1/status", NewRoute(api, baseChain, StatusGET))
	router.POST("/v1/import/sender", NewRoute(api, baseChain, ImportSenderPOST))
	router.GET("/v1/report/usage", NewRoute(api, baseChain, UsageReportGET))
	router.GET("/v1/report/usage/:account_id", NewRoute(api, baseChain, UsageReportGET))

	return api
}

// Route is some bullshit type so we can make httprouter chain with alice
// and still access params and such
type Route struct {
	params   httprouter.Params
	w        http.ResponseWriter
	r        *http.Request
	endpoint func(*Route)
	api      *AdminAPI
}

// NotFoundRoute route for convenience
func NotFoundRoute(r *Route) {
	http.Error(r.w, "Route not found", http.StatusNotFound)
}

// NewRoute returns a handler wrapped in middleware
// it's super important that route objects are created when requests come in
// if they are created before, the same object is shared for all requests
// and you fuck yourself with concurrent access... derp!
func NewRoute(api *AdminAPI, chain alice.Chain, endpoint func(*Route)) httprouter.Handle {
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
func NewPlainRoute(api *AdminAPI, chain alice.Chain, endpoint func(*Route)) http.HandlerFunc {
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
