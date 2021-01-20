package main

import (
	"log"

	accountRPC "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
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
	RPCPort            int    `envconfig:"RPC_PORT"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	RedisURL           string `envconfig:"REDIS_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	WebhookRPCHost     string `envconfig:"WEBHOOK_RPC_HOST"`
	WebhookRPCPort     int    `envconfig:"WEBHOOK_RPC_PORT"`
	SenderRPCHost      string `envconfig:"SENDER_RPC_HOST"`
	SenderRPCPort      int    `envconfig:"SENDER_RPC_PORT"`
	AccountRPCHost     string `envconfig:"ACCOUNT_RPC_HOST"`
	AccountRPCPort     int    `envconfig:"ACCOUNT_RPC_PORT"`
	TrackLinkDomain    string `envconfig:"TRACKLINK_DOMAIN"`
	OptOutLinkDomain   string `envconfig:"OPTOUTLINK_DOMAIN"`
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

	port := env.RPCPort

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	tracer, closer, err := jaeger.Connect(smsRPC.Name)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}
	defer closer.Close()

	wrpc := webhookpb.NewServiceClient(
		rpcbuilder.NewClientConn(env.WebhookRPCHost, env.WebhookRPCPort, tracer),
	)
	srpc := senderpb.NewServiceClient(
		rpcbuilder.NewClientConn(env.SenderRPCHost, env.SenderRPCPort, tracer),
	)
	arpc := accountRPC.New(env.AccountRPCHost, env.AccountRPCPort)
	tlrpc := tracklinkRPC.NewClient(env.TrackLinkRPCHost, env.TrackLinkRPCPort)
	orpc := optOutRPC.NewClient(env.OptOutRPCHost, env.OptOutRPCPort)

	features := smsRPC.SMSFeatures{
		TrackLinkDomain:  env.TrackLinkDomain,
		OptOutLinkDomain: env.OptOutLinkDomain,
	}

	smsrpc, err := smsRPC.NewService(features, env.PostgresURL, rabbitmq, wrpc, arpc, tlrpc, env.RedisURL, orpc, srpc)
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
