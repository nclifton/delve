package tracing

import (
	"fmt"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

func New(tracer opentracing.Tracer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return nethttp.Middleware(
			tracer,
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
