package main

import (
	"log"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/adminapi"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"

	"github.com/kelseyhightower/envconfig"
)

var gitref = "unset" // set with go linker in build script

type Env struct {
	Port string `envconfig:"PORT"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("ADMINAPI", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	app := adminapi.NewAdminAPI(&adminapi.AdminAPIOptions{
		NrApp: newrelicM,
	})

	log.Printf("%s service initialised and available on port %s", "adminapi", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, app.Handler()))
}
