package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/optout/inbound"
	rpc "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	OptoutRPCAddress  string `envconfig:"OPTOUT_RPC_ADDRESS"`
	OptoutInboundPort int    `envconfig:"INBOUND_PORT"`

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

	port := strconv.Itoa(env.OptoutInboundPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	opts := inbound.InboundOptions{
		OptOutClient: rpc.NewClient(env.OptoutRPCAddress),
		NrApp:        newrelicM,
	}

	server := inbound.NewInboundAPI(&opts)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", "sms inbound http", err)
	}

	log.Printf("%s service initialised and available on port %s", "sms inbound http", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
