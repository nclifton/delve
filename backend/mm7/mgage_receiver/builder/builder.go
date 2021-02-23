package builder

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rest"
	"github.com/burstsms/mtmo-tp/backend/lib/restbuilder"
	"github.com/burstsms/mtmo-tp/backend/lib/valid"
	"github.com/burstsms/mtmo-tp/backend/mm7/mgage_receiver/service"
	"github.com/kelseyhightower/envconfig"
)

func NewFromEnv() *serviceBuilder {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}
	return &serviceBuilder{conf: conf}
}

type Config struct {
}

type serviceBuilder struct {
	conf Config
}

func (b *serviceBuilder) Run(deps restbuilder.Deps) error {
	hb := rest.NewHandlerBuilder(&rest.HandlerConfig{
		Log:           deps.Log,
		Tracer:        deps.Tracer,
		JSONValidator: valid.Validate,
	})

	authRoute := hb().SetMiddleware(
		rest.NewTracingMiddleware(),
		rest.NewLoggingMiddleware(),
	).Handle
	baseRoute := hb().SetMiddleware(
		rest.NewTracingMiddleware(),
		rest.NewLoggingMiddleware(),
	).Handle

	service.Routes(deps.Router, baseRoute, authRoute)

	return nil
}
