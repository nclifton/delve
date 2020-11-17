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
	HTTPPort     int    `envconfig:"HTTP_PORT"`
	TemplatePath string `envconfig:"TEMPLATE_PATH"`
	DREndpoint   string `envconfig:"DR_ENDPOINT"`
	NRName       string `envconfig:"NR_NAME"`
	NRLicense    string `envconfig:"NR_LICENSE"`
	NRTracing    bool   `envconfig:"NR_TRACING"`
}

func main() {
	var env Env
	err := envconfig.Process("tualet", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := strconv.Itoa(env.HTTPPort)

	newrelicM := newrelicagent.New(&newrelicagent.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	opts := tualet.TeclooAPIOptions{
		NrApp:        newrelicM,
		TemplatePath: env.TemplatePath,
		DREndpoint:   env.DREndpoint,
	}

	server := t.NewTeclooAPI(&opts)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", "tualet api http", err)
	}

	log.Printf("%s service initialised and available on port %s", "tualet api http", port)
	log.Println("Tualet API: listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))

}
