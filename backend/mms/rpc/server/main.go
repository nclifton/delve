package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc"
	optOut "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
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
	TrackLinkRPCHost   string `envconfig:"TRACK_LINK_RPC_HOST"`
	TrackLinkRPCPort   int    `envconfig:"TRACK_LINK_RPC_PORT"`
	OptOutRPCHost      string `envconfig:"OPT_OUT_RPC_HOST"`
	OptOutRPCPort      int    `envconfig:"OPT_OUT_RPC_PORT"`
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
		Webhook:   webhook.NewClient(env.WebhookRPCHost, env.WebhookRPCPort),
		TrackLink: tracklink.NewClient(env.TrackLinkRPCHost, env.TrackLinkRPCPort),
		OptOut:    optOut.NewClient(env.OptOutRPCHost, env.OptOutRPCPort),
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
