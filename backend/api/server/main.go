package main

import (
	"log"
	"net/http"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/api"
	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"

	"github.com/kelseyhightower/envconfig"
)

const apiName = "REST API"

var gitref = "unset" // set with go linker in build script

type Env struct {
	APIPort           string `envconfig:"API_PORT"`
	AccountRPCAddress string `envconfig:"ACCOUNT_RPC_ADDRESS"`
	SMSRPCAddress     string `envconfig:"SMS_RPC_ADDRESS"`
	MMSRPCAddress     string `envconfig:"MMS_RPC_ADDRESS"`
	WebhookRPCAddress string `envconfig:"WEBHOOK_RPC_ADDRESS"`
	SenderRPCAddress  string `envconfig:"SENDER_RPC_ADDRESS"`
	NRName            string `envconfig:"NR_NAME"`
	NRLicense         string `envconfig:"NR_LICENSE"`
	NRTracing         bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("API", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	tracer, closer, err := jaeger.Connect(apiName)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", apiName, err)
	}
	defer closer.Close()

	app := api.New(&api.Options{
		Tracer:        tracer,
		Gitref:        gitref,
		AccountClient: account.New(env.AccountRPCAddress),
		SMSClient:     sms.New(env.SMSRPCAddress),
		MMSClient:     mms.New(env.MMSRPCAddress),
		WebhookClient: webhookpb.NewServiceClient(
			rpcbuilder.NewClientConn(env.WebhookRPCAddress, tracer),
		),
		SenderClient: senderpb.NewServiceClient(
			rpcbuilder.NewClientConn(env.SenderRPCAddress, tracer),
		),
		NrApp: newrelicM,
	})

	log.Printf("%s service initialised and available on port %s", "api", env.APIPort)
	log.Fatal(http.ListenAndServe(":"+env.APIPort, app.Handler()))
}
