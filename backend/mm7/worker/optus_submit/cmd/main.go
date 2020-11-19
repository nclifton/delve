package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/mm7/worker"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	mm7w "github.com/burstsms/mtmo-tp/backend/mm7/worker/optus_submit"
	"github.com/kelseyhightower/envconfig"
)

var Name = "mm7-worker-optus-submit"

type Env struct {
	RabbitURL             string `envconfig:"RABBIT_URL"`
	RabbitExchange        string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType    string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RabbitPrefetchedCount int    `envconfig:"RABBIT_PREFETCHED_COUNT"`
}

func main() {
	log.Printf("starting worker: %s", Name)

	var env Env
	err := envconfig.Process("mm7", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise rabbit worker: %s reason: %s\n", Name, err)
	}

	opts := rabbit.ConsumeOptions{
		PrefetchCount: env.RabbitPrefetchedCount,
		Exchange:      env.RabbitExchange,
		ExchangeType:  env.RabbitExchangeType,
		RouteKey:      worker.QueueNameSubmitOptus,
		QueueName:     worker.QueueNameSubmitOptus,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)
	w.Run(opts, mm7w.NewHandler(nil))
}
