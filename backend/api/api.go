package api

import (
	"net/http"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/api/middleware/auth"
	"github.com/burstsms/mtmo-tp/backend/api/middleware/context"
	"github.com/burstsms/mtmo-tp/backend/api/middleware/tracing"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/middleware/recovery"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go"
)

// Options will hold some state for our http handler
type Options struct {
	Tracer        opentracing.Tracer
	Gitref        string
	AccountClient *account.Client
	SMSClient     *sms.Client
	MMSClient     *mms.Client
	WebhookClient webhookpb.ServiceClient
	NrApp         func(http.Handler) http.Handler
}

type RPCClients struct {
	account *account.Client
	sms     *sms.Client
	mms     *mms.Client
	webhook webhookpb.ServiceClient
}

// API wraps an instance of our api app
type API struct {
	opts   *Options
	router *httprouter.Router
	log    *logger.StandardLogger
	RPCClients
}

// Handler exposes the router
func (a *API) Handler() http.Handler {
	return a.router
}

// New creates our api "app", i.e. the http handler
func New(opts *Options) *API {
	clients := RPCClients{
		account: opts.AccountClient,
		sms:     opts.SMSClient,
		mms:     opts.MMSClient,
		webhook: opts.WebhookClient,
	}

	api := &API{
		opts:       opts,
		log:        logger.NewLogger(),
		RPCClients: clients,
	}

	newrelicM := opts.NrApp

	tracingM := tracing.New(opts.Tracer)

	// maybe not needed with httprouter panic handler
	recoveryM := recovery.New(&recovery.Options{
		PrintStack: true,
	})

	authM := auth.New(&auth.Options{
		UseAPIKey: true,
		FindAccountByAPIKey: func(key string) (*account.Account, error) {
			reply, err := api.account.FindByAPIKey(key)
			return reply.Account, err
		},
	})

	// define the middleware chains for our api endpoints
	// we can group them and expand chains
	baseChain := alice.New(newrelicM, recoveryM, tracingM)
	defaultChain := baseChain.Append(context.ClearHandler)
	authChain := baseChain.Append(authM).Append(context.ClearHandler)

	// configure routes
	router := httprouter.New() // /v1/
	api.router = router

	// we need non trailing slash versions due to httprouter
	// catch-all matching and auto redirects
	router.Handle("OPTIONS", "/v1", NewRoute(api, defaultChain, EmptyRoute))
	router.Handle("OPTIONS", "/v1/*path", NewRoute(api, defaultChain, EmptyRoute))

	// we also need a route and chain for 404
	router.NotFound = NewPlainRoute(api, defaultChain, NotFoundRoute)

	// ------ routes without auth
	router.GET("/", NewRoute(api, defaultChain, IndexRoute))

	// ------ authenticated routes
	router.POST("/v1/sms", NewRoute(api, authChain, SMSPOST))

	router.POST("/v1/mms", NewRoute(api, authChain, MMSPOST))

	router.POST("/v1/webhook/test", NewRoute(api, authChain, TestPublishOptOutWebhookPOST))
	router.POST("/v1/webhook", NewRoute(api, authChain, WebhookCreatePOST))
	router.PUT("/v1/webhook/:id", NewRoute(api, authChain, WebhookUpdatePUT))
	router.GET("/v1/webhook/:id", NewRoute(api, authChain, WebhookGET))
	router.GET("/v1/webhook", NewRoute(api, authChain, WebhookListGET))
	router.DELETE("/v1/webhook/:id", NewRoute(api, authChain, WebhookDELETE))

	return api
}

// we have a function like this because over rpc you can't compare errors directly
/*func errCmp(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	return e1.Error() == e2.Error()
}*/
