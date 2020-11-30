package main

import (
	"log"

	mm7c "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	mmsc "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/mms/worker"
	mmsw "github.com/burstsms/mtmo-tp/backend/mms/worker/mms_send"
	"github.com/kelseyhightower/envconfig"
)

var Name = "mms-send"

type Env struct {
	RabbitURL             string `envconfig:"RABBIT_URL"`
	RabbitExchange        string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType    string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RabbitPrefetchedCount int    `envconfig:"RABBIT_PREFETCHED_COUNT"`
	MM7RPCHost            string `envconfig:"MM7_RPC_HOST"`
	MM7RPCPort            int    `envconfig:"MM7_RPC_PORT"`
	MMSRPCHost            string `envconfig:"RPC_HOST"`
	MMSRPCPort            int    `envconfig:"RPC_PORT"`
}

func main() {
	log.Printf("starting worker: %s", Name)

	var env Env
	err := envconfig.Process("mms", &env)
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
		RouteKey:      worker.MMSSendRouteKey,
		QueueName:     worker.MMSSendQueueName,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)

	mm7cli := mm7c.NewClient(env.MM7RPCHost, env.MM7RPCPort)
	mmscli := mmsc.New(env.MMSRPCHost, env.MMSRPCPort)

	w.Run(opts, mmsw.NewHandler(mm7cli, mmscli))
}
