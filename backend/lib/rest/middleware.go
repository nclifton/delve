package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

type Authenticator = func(key string) (interface{}, error)

func NewTracingMiddleware() MiddlewareConstructor {
	return func(hc *HandlerConfig) Middleware {
		return func(h http.Handler) http.Handler {
			return nethttp.Middleware(
				hc.Tracer,
				h,
				nethttp.OperationNameFunc(func(r *http.Request) string {
					return fmt.Sprintf("HTTP %s %s", r.Method, r.RequestURI)
				}),
				nethttp.MWSpanObserver(func(sp opentracing.Span, r *http.Request) {
					sp.SetTag("http.uri", r.URL.EscapedPath())
				}),
			)
		}
	}
}

func NewLoggingMiddleware() MiddlewareConstructor {
	return func(_ *HandlerConfig) Middleware {
		return func(h http.Handler) http.Handler {
			// TODO: use our own logger here somehow
			return handlers.LoggingHandler(os.Stdout, h)
		}
	}
}

func NewAuthMiddleware(authenticator Authenticator) MiddlewareConstructor {
	return func(hc *HandlerConfig) Middleware {
		return func(h http.Handler) http.Handler {
			return &authMiddleware{handler: h, hc: hc, authenticator: authenticator}
		}
	}
}

type authMiddleware struct {
	handler       http.Handler
	hc            *HandlerConfig
	authenticator Authenticator
}

func (m *authMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hctx := NewHandlerContext(w, r, m.hc)
	key := r.Header.Get("x-api-key")
	if key == "" {
		hctx.WriteJSONError("x-api-key header required", http.StatusBadRequest, nil)
		return
	}

	// verify the token by fetching account via api key
	reply, err := m.authenticator(key)
	if err != nil {
		hctx.WriteJSONError("Unauthorized", http.StatusUnauthorized, nil)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), AuthKey, reply))

	m.handler.ServeHTTP(w, r)
}
