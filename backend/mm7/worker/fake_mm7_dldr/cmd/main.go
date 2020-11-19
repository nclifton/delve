package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"

	"github.com/burstsms/mtmo-tp/backend/mm7/worker"

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
	RPCHost               string `envconfig:"RPC_HOST"`
	RPCPort               int    `envconfig:"RPC_PORT"`
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
		log.Fatalf("failed to initialise rabbit: %s reason: %s\n", Name, err)
	}

	opts := rabbit.ConsumeOptions{
		PrefetchCount: env.RabbitPrefetchedCount,
		Exchange:      env.RabbitExchange,
		ExchangeType:  env.RabbitExchangeType,
		RouteKey:      worker.QueueNameDLDRFake,
		QueueName:     worker.QueueNameDLDRFake,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)

	cli := client.NewClient(env.RPCHost, env.RPCPort)

	w.Run(opts, mm7w.NewHandler(cli))
}
