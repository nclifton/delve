package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/tualet"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	HTTPPort    int    `envconfig:"HTTP_PORT"`
	DLREndpoint string `envconfig:"DLR_ENDPOINT"`
	MOEndpoint  string `envconfig:"MO_ENDPOINT"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("tualet", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	port := strconv.Itoa(env.HTTPPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	opts := tualet.TualetAPIOptions{
		NrApp:       newrelicM,
		DLREndpoint: env.DLREndpoint,
		MOEndpoint:  env.MOEndpoint,
	}

	server := tualet.NewTualetAPI(&opts)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", "tualet", err)
	}

	log.Printf("%s service initialised and available on port %s", "tualet", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
