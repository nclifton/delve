package main

import (
	"log"

	agent "github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/processmo"
	"github.com/kelseyhightower/envconfig"
)

// TODO put this in config
var workerName = "processmo"

type Env struct {
	RabbitURL string `envconfig:"RABBIT_URL"`
	SMSHost   string `envconfig:"SMS_RPC_HOST"`
	SMSPort   int    `envconfig:"SMS_RPC_PORT"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("sms", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	// TODO service/worker level config for this url
	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise rabbit worker: %s reason: %s\n", workerName, err)
	}

	// TODO put this data in config
	opts := rabbit.ConsumeOptions{
		PrefetchCount: 1,
		QueueName:     msg.MOMessage.Queue,
		Exchange:      msg.MOMessage.Exchange,
		ExchangeType:  msg.MOMessage.ExchangeType,
		RouteKey:      msg.MOMessage.RouteKey,
		RetryScale:    rabbit.RetryScale,
	}

	nrOpts := &agent.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	}

	worker := rabbit.NewWorker(workerName, rabbitmq, nrOpts)

	client := rpc.New(env.SMSHost, env.SMSPort)

	log.Println("Service started")
	worker.Run(opts, processmo.NewHandler(client))
}
