package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/track_link/inbound"
	rpc "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	TrackLinkHost        string `envconfig:"RPC_HOST"`
	TrackLinkPort        int    `envconfig:"RPC_PORT"`
	TrackLinkInboundPort int    `envconfig:"INBOUND_PORT"`

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

	port := strconv.Itoa(env.TrackLinkInboundPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	opts := inbound.TrackLinkOptions{
		TrackLinkClient: rpc.NewClient(env.TrackLinkHost, env.TrackLinkPort),
		NrApp:           newrelicM,
	}

	server := inbound.NewTrackLinkAPI(&opts)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", "track link inbound http", err)
	}

	log.Printf("%s service initialised and available on port %s", "track link inbound http", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
