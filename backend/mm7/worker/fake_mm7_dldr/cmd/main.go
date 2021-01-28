package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"

	"github.com/burstsms/mtmo-tp/backend/mm7/worker"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	mm7w "github.com/burstsms/mtmo-tp/backend/mm7/worker/fake_mm7_dldr"
	"github.com/kelseyhightower/envconfig"
)

var Name = "fake-mm7-dldr-worker"

type Env struct {
	RabbitURL             string `envconfig:"RABBIT_URL"`
	RabbitExchange        string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType    string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RabbitPrefetchedCount int    `envconfig:"RABBIT_PREFETCHED_COUNT"`
	MM7RPCAddress         string `envconfig:"MM7_RPC_ADDRESS"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("mm7", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	// Register service with New Relic
	nr.CreateApp(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise rabbit: %s reason: %s\n", Name, err)
	}

	opts := rabbit.ConsumeOptions{
		PrefetchCount: env.RabbitPrefetchedCount,
		Exchange:      env.RabbitExchange,
		ExchangeType:  env.RabbitExchangeType,
		RouteKey:      worker.QueueNameDLDRFake,
		QueueName:     worker.QueueNameDLDRFake,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)

	cli := client.NewClient(env.MM7RPCAddress)

	log.Println("Service started")
	w.Run(opts, mm7w.NewHandler(cli))
}
