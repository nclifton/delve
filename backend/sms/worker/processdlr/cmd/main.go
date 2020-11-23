package main

import (
	"log"

	agent "github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/processdlr"
	"github.com/kelseyhightower/envconfig"
)

// TODO put this in config
var workerName = "processdlr"

type Env struct {
	RabbitURL string `envconfig:"RABBIT_URL"`
	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
	SMSHost   string `envconfig:"SMS_HOST"`
	SMSPort   int    `envconfig:"SMS_RPC_PORT"`
}

func main() {
	var env Env
	err := envconfig.Process("sms", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	log.Printf("starting worker: %s", workerName)

	// TODO service/worker level config for this url
	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise rabbit worker: %s reason: %s\n", workerName, err)
	}

	// TODO put this data in config
	opts := rabbit.ConsumeOptions{
		PrefetchCount: 1,
		QueueName:     msg.DLRMessage.Queue,
		Exchange:      msg.DLRMessage.Exchange,
		ExchangeType:  msg.DLRMessage.ExchangeType,
		RouteKey:      msg.DLRMessage.RouteKey,
		RetryScale:    rabbit.RetryScale,
	}

	nrOpts := &agent.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	}

	worker := rabbit.NewWorker(workerName, rabbitmq, nrOpts)

	client := rpc.New(env.SMSHost, env.SMSPort)

	worker.Run(opts, processdlr.NewHandler(client))
}
