package run

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/servicebuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/queue"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/kelseyhightower/envconfig"
)

type webhookEnv struct {
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
}

func Server(deps servicebuilder.Deps) error {
	var env webhookEnv
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	rqueue := queue.NewRabbitQueue(deps.RabbitConn, rabbit.PublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
		Tracer:       deps.Tracer,
	})
	pdb := db.NewSQLDB(deps.PostgresConn)

	webhookpb.RegisterServiceServer(deps.Server, service.NewWebhookService(pdb, rqueue))

	return nil
}
