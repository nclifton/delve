package builder

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/queue"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostRabbitExchange     string `envconfig:"POST_RABBIT_EXCHANGE"`
	PostRabbitExchangeType string `envconfig:"POST_RABBIT_EXCHANGE_TYPE"`
}

type builder struct {
	conf Config
}

func NewService(config Config) *builder {
	return &builder{conf: config}
}

func NewServiceFromEnv() *builder {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return NewService(config)
}

func (b *builder) Run(deps rpcbuilder.Deps) error {
	rqueue := queue.NewRabbitQueue(deps.RabbitConn, rabbit.PublishOptions{
		Exchange:     b.conf.PostRabbitExchange,
		ExchangeType: b.conf.PostRabbitExchangeType,
		Tracer:       deps.Tracer,
	})
	pdb := db.NewSQLDB(deps.PostgresConn)

	webhookpb.RegisterServiceServer(deps.Server, service.NewWebhookService(pdb, rqueue))

	return nil
}
