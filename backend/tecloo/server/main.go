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

	Name    string `envconfig:"NAME"`
	License string `envconfig:"LICENSE"`
	Tracing bool   `envconfig:"TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("tecloo", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	port := strconv.Itoa(env.HTTPPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.Name,
		NewRelicLicense:          env.License,
		DistributedTracerEnabled: env.Tracing,
	})

	opts := tecloo.TeclooAPIOptions{
		NrApp:        newrelicM,
		TemplatePath: env.TemplatePath,
		DREndpoint:   env.DREndpoint,
	}

	server := tecloo.NewTeclooAPI(&opts)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", "tecloo", err)
	}

	log.Printf("%s service initialised and available on port %s", "tecloo", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
