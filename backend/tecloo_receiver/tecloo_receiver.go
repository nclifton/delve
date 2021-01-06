package tecloo_receiver

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/recovery"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type TeclooReceiverAPIOptions struct {
	NrApp        func(http.Handler) http.Handler
	TemplatePath string
	Rabbit       rabbit.Conn
	RabbitOpts   RabbitPublishOptions
}

type TeclooReceiverTemplates struct {
	DeliverResponse        *template.Template
	DeliveryReportResponse *template.Template
}

// API wraps an instance of our api app
type TeclooReceiverAPI struct {
	opts      *TeclooReceiverAPIOptions
	router    *httprouter.Router
	log       *logger.StandardLogger
	templates *TeclooReceiverTemplates
}

type RabbitPublishOptions = rabbit.PublishOptions

func (a *TeclooReceiverAPI) Publish(msg interface{}, header map[string]interface{}, dontEncodeJson bool, routeKey string) error {
	publishOpts := RabbitPublishOptions{
		Exchange:       a.opts.RabbitOpts.Exchange,
		ExchangeType:   a.opts.RabbitOpts.ExchangeType,
		Headers:        header,
		RouteKey:       routeKey,
		DontEncodeJson: dontEncodeJson,
	}

	return rabbit.Publish(a.opts.Rabbit, publishOpts, msg)
}

// Handler exposes the router
func (a *TeclooReceiverAPI) Handler() http.Handler {
	return a.router
}

// New creates our api "app", i.e. the http handler
func NewTeclooReceiverAPI(opts *TeclooReceiverAPIOptions) *TeclooReceiverAPI {
	templates := &TeclooReceiverTemplates{
		DeliverResponse:        template.Must(template.ParseFiles(fmt.Sprintf(`%s/mm7_deliver_response.soap.tmpl`, opts.TemplatePath))),
		DeliveryReportResponse: template.Must(template.ParseFiles(fmt.Sprintf(`%s/mm7_delivery_report_response.soap.tmpl`, opts.TemplatePath))),
	}

	api := &TeclooReceiverAPI{
		opts:      opts,
		log:       logger.NewLogger(),
		templates: templates,
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
	router.POST("/v1/mms/inbound", NewRoute(api, baseChain, InboundPOST))

	return api
}

// Route is some bullshit type so we can make httprouter chain with alice
// and still access params and such
type Route struct {
	params   httprouter.Params
	w        http.ResponseWriter
	r        *http.Request
	endpoint func(*Route)
	api      *TeclooReceiverAPI
}

// NotFoundRoute route for convenience
func NotFoundRoute(r *Route) {
	http.Error(r.w, "Route not found", http.StatusNotFound)
}

// NewRoute returns a handler wrapped in middleware
// it's super important that route objects are created when requests come in
// if they are created before, the same object is shared for all requests
// and you fuck yourself with concurrent access... derp!
func NewRoute(api *TeclooReceiverAPI, chain alice.Chain, endpoint func(*Route)) httprouter.Handle {
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
func NewPlainRoute(api *TeclooReceiverAPI, chain alice.Chain, endpoint func(*Route)) http.HandlerFunc {
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
