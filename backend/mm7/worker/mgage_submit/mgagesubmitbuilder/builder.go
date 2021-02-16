package mgagesubmitbuilder

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker/mgage_submit/handler"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
}

type service struct {
	conf Config
}

func NewBuilderFromEnv() *service {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return New(config)
}

func New(config Config) *service {
	return &service{conf: config}
}

func (s *service) Run(deps workerbuilder.Deps) error {

	handler := handler.New()

	deps.Health.SetServiceReady(true)
	deps.Worker.Run(deps.ConsumeOptions, handler)

	return nil
}
