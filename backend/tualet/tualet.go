package tualet

import (
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/middleware/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/recovery"
	belogger "github.com/burstsms/mtmo-tp/backend/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type TualetAPIOptions struct {
	NrApp       func(http.Handler) http.Handler
	DLREndpoint string
	MOEndpoint  string
	Client      *http.Client
}

// API wraps an instance of our api app
type TualetAPI struct {
	opts   *TualetAPIOptions
	router *httprouter.Router
	log    *belogger.StandardLogger
	client *http.Client
}

// Handler exposes the router
func (a *TualetAPI) Handler() http.Handler {
	return a.router
}

// New creates our api "app", i.e. the http handler
func NewTualetAPI(opts *TualetAPIOptions) *TualetAPI {

	client := opts.Client
	if client == nil {
		client = http.DefaultClient
	}

	api := &TualetAPI{
		opts:   opts,
		log:    belogger.NewLogger(),
		client: client,
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

	if opts.NrApp != nil {
		baseChain.Append(opts.NrApp)
	}

	// configure routes
	router := httprouter.New() // /v1/
	api.router = router

	// we also need a route and chain for 404
	router.NotFound = NewPlainRoute(api, baseChain, NotFoundRoute)
	router.GET("/v1/status", NewRoute(api, baseChain, StatusGET))
	router.GET("/v1/fakemo", NewRoute(api, baseChain, HandsetGET))
	router.GET("/api", NewRoute(api, baseChain, SubmitGET))

	return api
}

// Route is some bullshit type so we can make httprouter chain with alice
// and still access params and such
type Route struct {
	params   httprouter.Params
	w        http.ResponseWriter
	r        *http.Request
	endpoint func(*Route)
	api      *TualetAPI
}

// NotFoundRoute route for convenience
func NotFoundRoute(r *Route) {
	http.Error(r.w, "Route not found", http.StatusNotFound)
}

// NewRoute returns a handler wrapped in middleware
// it's super important that route objects are created when requests come in
// if they are created before, the same object is shared for all requests
// and you fuck yourself with concurrent access... derp!
func NewRoute(api *TualetAPI, chain alice.Chain, endpoint func(*Route)) httprouter.Handle {
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
func NewPlainRoute(api *TualetAPI, chain alice.Chain, endpoint func(*Route)) http.HandlerFunc {
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
