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
	NRName            string `envconfig:"NR_NAME"`
	NRLicense         string `envconfig:"NR_LICENSE"`
	NRTracing         bool   `envconfig:"NR_TRACING"`
	OptOutHost        string `envconfig:"RPC_HOST"`
	OptOutPort        int    `envconfig:"RPC_PORT"`
	OptOutInboundPort int    `envconfig:"INBOUND_PORT"`
}

func main() {
	var env Env
	err := envconfig.Process("optout", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := strconv.Itoa(env.OptOutInboundPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	opts := inbound.InboundOptions{
		OptOutClient: rpc.New(env.OptOutHost, env.OptOutPort),
		NrApp:        newrelicM,
	}

	server := inbound.NewInboundAPI(&opts)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", "sms inbound http", err)
	}

	log.Printf("%s service initialised and available on port %s", "sms inbound http", port)
	log.Println("SMS Inbound HTTP API: listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
