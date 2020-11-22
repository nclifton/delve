package main

import (
	"log"
	"net/http"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/api"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"

	"github.com/kelseyhightower/envconfig"
)

var gitref = "unset" // set with go linker in build script

type Env struct {
	Port string `envconfig:"PORT"`

	AccountHost string `envconfig:"ACCOUNT_HOST"`
	AccountPort int    `envconfig:"ACCOUNT_PORT"`

	SMSHost string `envconfig:"SMS_HOST"`
	SMSPort int    `envconfig:"SMS_PORT"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	var env Env
	err := envconfig.Process("API", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	app := api.New(&api.Options{
		Gitref:        gitref,
		AccountClient: account.New(env.AccountHost, env.AccountPort),
		SMSClient:     sms.New(env.SMSHost, env.SMSPort),
		NrApp:         newrelicM,
	})

	log.Println("API: listening on", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, app.Handler()))
}
