package rest

import (
	"context"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
)

const (
	ParamsKey = "params"
	AuthKey   = "auth"
	LogKey    = "logger"
)

type Handle = func(h Handler) httprouter.Handle
type Handler = func(rc *HandlerContext)
type Middleware func(http.Handler) http.Handler
type MiddlewareConstructor func(*HandlerConfig) Middleware

// HandlerConfig global config injected into new handler chains
type HandlerConfig struct {
	Log           *logger.StandardLogger
	Tracer        opentracing.Tracer
	JSONValidator Validator
}

// NewHandlerBuilder returns a generator function for building new handlers
func NewHandlerBuilder(config *HandlerConfig) func() *handler {
	return func() *handler {
		return &handler{config: config}
	}
}

type handler struct {
	middleware []Middleware
	config     *HandlerConfig
}

// SetMiddleware sets middleware to be chained for route
func (hr *handler) SetMiddleware(ms ...MiddlewareConstructor) *handler {
	for _, m := range ms {
		hr.middleware = append(hr.middleware, m(hr.config))
	}
	return hr
}

// Handle chains previously set middleware and outputs a httprouter.Handle for router
func (hr *handler) Handle(h Handler) httprouter.Handle {
	rt := HandlerWrapper(h, hr.config)
	for _, v := range hr.middleware {
		if v != nil {
			rt = v(rt)
		}
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r = r.WithContext(context.WithValue(r.Context(), ParamsKey, ps))

		rt.ServeHTTP(w, r)
	}
}

// HandlerWrapper wraps custom handlers with http.Handler for middleware chaining
func HandlerWrapper(h Handler, config *HandlerConfig) http.Handler {
	return handlerWrapper{h: h, config: config}
}

type handlerWrapper struct {
	h      Handler
	config *HandlerConfig
}

func (rt handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.h(NewHandlerContext(w, r, rt.config))
}

// HTTPHandlerWrapper util for registering http.Handler routes
func HTTPHandlerWrapper(h httprouter.Handle) http.Handler {
	return httpHandlerWrapper{h: h}
}

type httpHandlerWrapper struct {
	h httprouter.Handle
}

func (hw httpHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps, ok := r.Context().Value(ParamsKey).(httprouter.Params)
	if !ok {
		ps = nil
	}

	hw.h(w, r, ps)
}
