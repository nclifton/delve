package run

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/servicebuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	"github.com/kelseyhightower/envconfig"
)

type senderEnv struct {
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
}

func Server(deps servicebuilder.Deps) error {
	var env senderEnv
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	pdb := db.NewSQLDB(deps.PostgresConn)

	senderpb.RegisterServiceServer(deps.Server, service.NewSenderService(pdb))

	return nil
}
