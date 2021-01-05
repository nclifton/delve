package jaeger

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func Connect(serviceName string) (opentracing.Tracer, io.Closer, error) {
	// metricsFactory := prometheus.New()
	tracer, closer, err := config.Configuration{
		ServiceName: serviceName,
	}.NewTracer(
	// config.Metrics(metricsFactory),
	)

	return tracer, closer, err
}
