package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc"
	optOut "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"

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
	SenderRPCHost      string `envconfig:"SENDER_RPC_HOST"`
	SenderRPCPort      int    `envconfig:"SENDER_RPC_PORT"`
	TrackLinkRPCHost   string `envconfig:"TRACK_LINK_RPC_HOST"`
	TrackLinkRPCPort   int    `envconfig:"TRACK_LINK_RPC_PORT"`
	OptOutRPCHost      string `envconfig:"OPT_OUT_RPC_HOST"`
	OptOutRPCPort      int    `envconfig:"OPT_OUT_RPC_PORT"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("mms", &env)
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

	port := env.RPCPort

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	rabbitOpts := rabbit.PublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
	}

	tracer, closer, err := jaeger.Connect(mmsRPC.Name)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}
	defer closer.Close()

	svc := mmsRPC.ConfigSvc{
		Webhook: webhookpb.NewServiceClient(
			rpcbuilder.NewClientConn(env.WebhookRPCHost, env.WebhookRPCPort, tracer),
		),
		Sender: senderpb.NewServiceClient(
			rpcbuilder.NewClientConn(env.SenderRPCHost, env.SenderRPCPort, tracer),
		),
		TrackLink: tracklink.NewClient(env.TrackLinkRPCHost, env.TrackLinkRPCPort),
		OptOut:    optOut.NewClient(env.OptOutRPCHost, env.OptOutRPCPort),
	}

	mmsrpc, err := mmsRPC.NewService(env.PostgresURL, rabbitmq, rabbitOpts, svc)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	server, err := rpc.NewServer(mmsrpc, port)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", mmsRPC.Name, port)
	server.Listen()
}
