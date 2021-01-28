package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/sms/inbound"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	SMSRPCAddress  string `envconfig:"SMS_RPC_ADDRESS"`
	SMSInboundPort int    `envconfig:"SMS_INBOUND_PORT"`

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

	port := strconv.Itoa(env.SMSInboundPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	opts := inbound.InboundOptions{
		SMSClient: rpc.New(env.SMSRPCAddress),
		NrApp:     newrelicM,
	}

	server := inbound.NewInboundAPI(&opts)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", "sms inbound http", err)
	}

	log.Printf("%s service initialised and available on port %s", "sms inbound http", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
