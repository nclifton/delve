package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	ooRPC "github.com/burstsms/mtmo-tp/backend/optout/rpc"
	webhookRPC "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort            int    `envconfig:"RPC_PORT"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	WebhookRPCHost     string `envconfig:"WEBHOOK_RPC_HOST"`
	WebhookRPCPort     int    `envconfig:"WEBHOOK_RPC_PORT"`
	TrackHost          string `envconfig:"TRACK_HOST"`
}

func main() {
	var env Env
	err := envconfig.Process("optout", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	wrpc := webhookRPC.NewClient(env.WebhookRPCHost, env.WebhookRPCPort)

	orpc, err := ooRPC.NewService(env.PostgresURL, env.TrackHost, wrpc)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", ooRPC.Name, err)
	}

	server, err := rpc.NewServer(orpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", ooRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", ooRPC.Name, port)
	server.Listen()
}
