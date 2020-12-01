package main

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/webhook/db"
	webhookRPC "github.com/burstsms/mtmo-tp/backend/webhook/rpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RabbitURL          string `envconfig:"RABBIT_URL"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RPCPort            int    `envconfig:"RPC_PORT"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("webhook", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	// Register service with New Relic
	nr.CreateApp(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	postgres, err := pgxpool.Connect(context.Background(), env.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", webhookRPC.Name, err)
	}

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", webhookRPC.Name, err)
	}

	webhookDB := db.New(postgres, rabbitmq)

	server, err := rpc.NewServer(webhookRPC.NewService(webhookDB), env.RPCPort)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", webhookRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", webhookRPC.Name, env.RPCPort)
	server.Listen()
}
