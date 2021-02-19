package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	optOutRPC "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	smsRPC "github.com/burstsms/mtmo-tp/backend/sms/rpc"
	tracklinkRPC "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	ContainerName       string `envconfig:"CONTAINER_NAME"`
	ContainerPort       int    `envconfig:"CONTAINER_PORT"`
	PostgresURL         string `envconfig:"POSTGRES_URL"`
	RabbitURL           string `envconfig:"RABBIT_URL"`
	RedisURL            string `envconfig:"REDIS_URL"`
	RabbitExchange      string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType  string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	WebhookRPCAddress   string `envconfig:"WEBHOOK_RPC_ADDRESS"`
	SenderRPCAddress    string `envconfig:"SENDER_RPC_ADDRESS"`
	TrackLinkDomain     string `envconfig:"TRACKLINK_DOMAIN"`
	OptOutLinkDomain    string `envconfig:"OPTOUTLINK_DOMAIN"`
	TrackLinkRPCAddress string `envconfig:"TRACK_LINK_RPC_ADDRESS"`
	OptOutRPCAddress    string `envconfig:"OPTOUT_RPC_ADDRESS"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("sms", &env)
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

	port := env.ContainerPort

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	tracer, closer, err := jaeger.Connect(env.ContainerName)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}
	defer closer.Close()

	wrpc := webhookpb.NewServiceClient(
		rpcbuilder.NewClientConn(env.WebhookRPCAddress, tracer),
	)
	srpc := senderpb.NewServiceClient(
		rpcbuilder.NewClientConn(env.SenderRPCAddress, tracer),
	)

	tlrpc := tracklinkRPC.NewClient(env.TrackLinkRPCAddress)
	orpc := optOutRPC.NewClient(env.OptOutRPCAddress)

	features := smsRPC.SMSFeatures{
		TrackLinkDomain:  env.TrackLinkDomain,
		OptOutLinkDomain: env.OptOutLinkDomain,
	}

	smsrpc, err := smsRPC.NewService(features, env.PostgresURL, rabbitmq, wrpc, tlrpc, env.RedisURL, orpc, srpc)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	server, err := rpc.NewServer(smsrpc, port)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", smsRPC.Name, port)
	server.Listen()
}
