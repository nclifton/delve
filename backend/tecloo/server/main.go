package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/tecloo"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	HTTPPort     int    `envconfig:"HTTP_PORT"`
	TemplatePath string `envconfig:"TEMPLATE_PATH"`
	DREndpoint   string `envconfig:"DR_ENDPOINT"`
}

type NREnv struct {
	Name    string `envconfig:"NAME"`
	License string `envconfig:"LICENSE"`
	Tracing bool   `envconfig:"TRACING"`
}

func main() {
	var env Env
	err := envconfig.Process("tecloo", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}
	var nrenv NREnv
	err = envconfig.Process("nr", &nrenv)
	if err != nil {
		log.Fatal("failed to read new relic env vars:", err)
	}

	port := strconv.Itoa(env.HTTPPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  nrenv.Name,
		NewRelicLicense:          nrenv.License,
		DistributedTracerEnabled: nrenv.Tracing,
	})

	opts := tecloo.TeclooAPIOptions{
		NrApp:        newrelicM,
		TemplatePath: env.TemplatePath,
		DREndpoint:   env.DREndpoint,
	}

	server := tecloo.NewTeclooAPI(&opts)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", "tecloo api http", err)
	}

	log.Printf("%s service initialised and available on port %s", "tecloo api http", port)
	log.Println("Tecloo API: listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))

}
