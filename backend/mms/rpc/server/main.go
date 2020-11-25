package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort            int    `envconfig:"RPC_PORT"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	WebhookRPCHost     string `envconfig:"WEBHOOK_RPC_HOST"`
	WebhookRPCPort     int    `envconfig:"WEBHOOK_RPC_PORT"`
}

func main() {
	var env Env
	err := envconfig.Process("mms", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	rabbitOpts := rabbit.PublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
	}

	svc := mmsRPC.ConfigSvc{
		Webhook: webhook.NewClient(env.WebhookRPCHost, env.WebhookRPCPort),
	}

	mmsrpc, err := mmsRPC.NewService(env.PostgresURL, rabbitmq, rabbitOpts, svc)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	server, err := rpc.NewServer(mmsrpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", mmsRPC.Name, port)
	server.Listen()
}
