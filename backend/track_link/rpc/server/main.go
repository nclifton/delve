package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	smsRPC "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	tlrpc "github.com/burstsms/mtmo-tp/backend/track_link/rpc"
	webhookRPC "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort        int    `envconfig:"RPC_PORT"`
	PostgresURL    string `envconfig:"POSTGRES_URL"`
	TrackHost      string `envconfig:"TRACK_HOST"`
	MMSRPCHost     string `envconfig:"MMS_RPC_HOST"`
	MMSRPCPort     int    `envconfig:"MMS_RPC_PORT"`
	SMSRPCHost     string `envconfig:"SMS_RPC_HOST"`
	SMSRPCPort     int    `envconfig:"SMS_RPC_PORT"`
	WebhookRPCHost string `envconfig:"WEBHOOK_RPC_HOST"`
	WebhookRPCPort int    `envconfig:"WEBHOOK_RPC_PORT"`
}

func main() {
	var env Env
	err := envconfig.Process("track_link", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	mmsrpc := mmsRPC.New(env.MMSRPCHost, env.MMSRPCPort)
	smsrpc := smsRPC.New(env.SMSRPCHost, env.SMSRPCPort)
	wrpc := webhookRPC.NewClient(env.WebhookRPCHost, env.WebhookRPCPort)

	srpc, err := tlrpc.NewService(env.PostgresURL, env.TrackHost, mmsrpc, smsrpc, wrpc)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", tlrpc.Name, err)
	}

	server, err := rpc.NewServer(srpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", tlrpc.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", tlrpc.Name, port)
	server.Listen()
}
