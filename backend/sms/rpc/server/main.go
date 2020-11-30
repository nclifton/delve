package main

import (
	"log"

	accountRPC "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	optOutRPC "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	smsRPC "github.com/burstsms/mtmo-tp/backend/sms/rpc"
	tracklinkRPC "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	webhookRPC "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"

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
	AccountRPCHost     string `envconfig:"ACCOUNT_RPC_HOST"`
	AccountRPCPort     int    `envconfig:"ACCOUNT_RPC_PORT"`
	TrackLinkDomain    string `envconfig:"TRACKLINK_DOMAIN"`
	OptOutLinkDomain   string `envconfig:"OPTOUTLINK_DOMAIN"`
	TrackLinkRPCHost   string `envconfig:"TRACK_LINK_RPC_HOST"`
	TrackLinkRPCPort   int    `envconfig:"TRACK_LINK_RPC_PORT"`
	OptOutRPCHost      string `envconfig:"OPT_OUT_RPC_HOST"`
	OptOutRPCPort      int    `envconfig:"OPT_OUT_RPC_PORT"`
}

func main() {
	var env Env
	err := envconfig.Process("sms", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	wrpc := webhookRPC.NewClient(env.WebhookRPCHost, env.WebhookRPCPort)
	arpc := accountRPC.New(env.AccountRPCHost, env.AccountRPCPort)
	tlrpc := tracklinkRPC.NewClient(env.TrackLinkRPCHost, env.TrackLinkRPCPort)
	orpc := optOutRPC.NewClient(env.OptOutRPCHost, env.OptOutRPCPort)

	features := smsRPC.SMSFeatures{
		TrackLinkDomain:  env.TrackLinkDomain,
		OptOutLinkDomain: env.OptOutLinkDomain,
	}

	srpc, err := smsRPC.NewService(features, env.PostgresURL, rabbitmq, wrpc, arpc, tlrpc, env.RedisURL, orpc)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	server, err := rpc.NewServer(srpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", smsRPC.Name, port)
	server.Listen()
}
