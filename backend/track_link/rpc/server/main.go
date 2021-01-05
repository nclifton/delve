package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
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
	TrackDomain    string `envconfig:"TRACKLINK_DOMAIN"`
	MMSRPCHost     string `envconfig:"MMS_RPC_HOST"`
	MMSRPCPort     int    `envconfig:"MMS_RPC_PORT"`
	SMSRPCHost     string `envconfig:"SMS_RPC_HOST"`
	SMSRPCPort     int    `envconfig:"SMS_RPC_PORT"`
	WebhookRPCHost string `envconfig:"WEBHOOK_RPC_HOST"`
	WebhookRPCPort int    `envconfig:"WEBHOOK_RPC_PORT"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("track_link", &env)
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

	tracer, closer, err := jaeger.Connect(tlrpc.Name)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", tlrpc.Name, err)
	}
	defer closer.Close()

	mmsrpc := mmsRPC.New(env.MMSRPCHost, env.MMSRPCPort)
	smsrpc := smsRPC.New(env.SMSRPCHost, env.SMSRPCPort)
	wrpc := webhookRPC.NewClient(env.WebhookRPCHost, env.WebhookRPCPort, tracer)

	srpc, err := tlrpc.NewService(env.PostgresURL, env.TrackDomain, mmsrpc, smsrpc, wrpc)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", tlrpc.Name, err)
	}

	server, err := rpc.NewServer(srpc, port)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", tlrpc.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", tlrpc.Name, port)
	server.Listen()
}
