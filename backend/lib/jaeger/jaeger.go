package jaeger

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func Connect(serviceName string) (opentracing.Tracer, io.Closer, error) {
	// metricsFactory := prometheus.New()

	conf, err := config.FromEnv()
	if err != nil {
		return nil, nil, err
	}
	conf.ServiceName = serviceName

	tracer, closer, err := conf.NewTracer(
	// config.Metrics(metricsFactory),
	)
	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, err
}
