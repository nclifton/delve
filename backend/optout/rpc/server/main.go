package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/lib/servicebuilder"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	ooRPC "github.com/burstsms/mtmo-tp/backend/optout/rpc"
	smsRPC "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort        int    `envconfig:"RPC_PORT"`
	PostgresURL    string `envconfig:"POSTGRES_URL"`
	WebhookRPCHost string `envconfig:"WEBHOOK_RPC_HOST"`
	WebhookRPCPort int    `envconfig:"WEBHOOK_RPC_PORT"`
	SMSRPCHost     string `envconfig:"SMS_RPC_HOST"`
	SMSRPCPort     int    `envconfig:"SMS_RPC_PORT"`
	MMSRPCHost     string `envconfig:"MMS_RPC_HOST"`
	MMSRPCPort     int    `envconfig:"MMS_RPC_PORT"`
	OptOutDomain   string `envconfig:"OPTOUTLINK_DOMAIN"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("optout", &env)
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

	tracer, closer, err := jaeger.Connect(ooRPC.Name)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", ooRPC.Name, err)
	}
	defer closer.Close()

	wrpc := webhookpb.NewServiceClient(
		servicebuilder.NewClientConn(env.WebhookRPCHost, env.WebhookRPCPort, tracer),
	)
	smsrpc := smsRPC.New(env.SMSRPCHost, env.SMSRPCPort)
	mmsrpc := mmsRPC.New(env.MMSRPCHost, env.MMSRPCPort)

	orpc, err := ooRPC.NewService(env.PostgresURL, env.OptOutDomain, wrpc, smsrpc, mmsrpc)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", ooRPC.Name, err)
	}

	server, err := rpc.NewServer(orpc, port)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", ooRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", ooRPC.Name, port)
	server.Listen()
}
