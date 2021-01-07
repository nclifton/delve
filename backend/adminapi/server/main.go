package main

import (
	"log"
	"net/http"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/adminapi"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"

	"github.com/kelseyhightower/envconfig"
)

var gitref = "unset" // set with go linker in build script

type Env struct {
	Port string `envconfig:"ADMINAPI_PORT"`

	AccountHost string `envconfig:"ACCOUNT_HOST"`
	AccountPort int    `envconfig:"ACCOUNT_PORT"`

	SMSHost string `envconfig:"SMS_HOST"`
	SMSPort int    `envconfig:"SMS_PORT"`

	MMSHost string `envconfig:"MMS_HOST"`
	MMSPort int    `envconfig:"MMS_PORT"`

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
		NrApp:         newrelicM,
		SMSClient:     sms.New(env.SMSHost, env.SMSPort),
		MMSClient:     mms.New(env.MMSHost, env.MMSPort),
		AccountClient: account.New(env.AccountHost, env.AccountPort),
	})

	log.Printf("%s service initialised and available on port %s", "adminapi", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, app.Handler()))
}
