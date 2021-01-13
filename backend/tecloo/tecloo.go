package tecloo

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/recovery"

	"github.com/justinas/alice"
)

type TeclooAPIOptions struct {
	NrApp        func(http.Handler) http.Handler
	TemplatePath string
	DREndpoint   string
	Client       *http.Client
}

type TeclooTemplates struct {
	SubmitResponse     *template.Template
	SubmitRequest      *template.Template
	SendDeliveryReport *template.Template
	SendDelivery       *template.Template
}

// API wraps an instance of our api app
type TeclooAPI struct {
	opts      *TeclooAPIOptions
	router    *httprouter.Router
	log       *logger.StandardLogger
	templates *TeclooTemplates
	client    *http.Client
}

// Handler exposes the router
func (a *TeclooAPI) Handler() http.Handler {
	return a.router
}

// New creates our api "app", i.e. the http handler
func NewTeclooAPI(opts *TeclooAPIOptions) *TeclooAPI {

	templates := &TeclooTemplates{
		SubmitResponse:     template.Must(template.ParseFiles(fmt.Sprintf(`%s/tecloo_submit_response.soap.tmpl`, opts.TemplatePath))),
		SubmitRequest:      template.Must(template.ParseFiles(fmt.Sprintf(`%s/tecloo_submit.soap.tmpl`, opts.TemplatePath))),
		SendDeliveryReport: template.Must(template.ParseFiles(fmt.Sprintf(`%s/tecloo_submit_dr.soap.tmpl`, opts.TemplatePath))),
		SendDelivery:       template.Must(template.ParseFiles(fmt.Sprintf(`%s/tecloo_delivery.soap.tmpl`, opts.TemplatePath))),
	}

	client := opts.Client
	if client == nil {
		client = http.DefaultClient
	}

	api := &TeclooAPI{
		opts:      opts,
		log:       logger.NewLogger(),
		templates: templates,
		client:    client,
	}

	// maybe not needed with httprouter panic handler
	recoveryM := recovery.New(&recovery.Options{
		PrintStack: true,
	})

	// define the middleware chains for our api endpoints
	// we can group them and expand chains
	baseChain := alice.New(recoveryM)

	if opts.NrApp != nil {
		baseChain.Append(opts.NrApp)
	}

	// configure routes
	router := httprouter.New() // /v1/
	api.router = router

	// we also need a route and chain for 404
	router.NotFound = NewPlainRoute(api, baseChain, NotFoundRoute)
	router.GET("/v1/status", NewRoute(api, baseChain, StatusGET))
	router.POST("/v1/mm7", NewRoute(api, baseChain, SubmitPOST))
	router.POST("/v1/handset/mms", NewRoute(api, baseChain, HandsetPOST))

	return api
}

// Route is some bullshit type so we can make httprouter chain with alice
// and still access params and such
type Route struct {
	params   httprouter.Params
	w        http.ResponseWriter
	r        *http.Request
	endpoint func(*Route)
	api      *TeclooAPI
}

// NotFoundRoute route for convenience
func NotFoundRoute(r *Route) {
	http.Error(r.w, "Route not found", http.StatusNotFound)
}

// NewRoute returns a handler wrapped in middleware
// it's super important that route objects are created when requests come in
// if they are created before, the same object is shared for all requests
// and you fuck yourself with concurrent access... derp!
func NewRoute(api *TeclooAPI, chain alice.Chain, endpoint func(*Route)) httprouter.Handle {
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
func NewPlainRoute(api *TeclooAPI, chain alice.Chain, endpoint func(*Route)) http.HandlerFunc {
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
