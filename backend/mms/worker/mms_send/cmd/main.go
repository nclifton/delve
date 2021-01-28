package main

import (
	"log"

	mm7c "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	mmsc "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
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
	MM7RPCAddress         string `envconfig:"MM7_RPC_ADDRESS"`
	MMSRPCAddress         string `envconfig:"MMS_RPC_ADDRESS"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("mms", &env)
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
		RouteKey:      worker.MMSSendRouteKey,
		QueueName:     worker.MMSSendQueueName,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)

	mm7cli := mm7c.NewClient(env.MM7RPCAddress)
	mmscli := mmsc.New(env.MMSRPCAddress)

	log.Println("Service started")
	w.Run(opts, mmsw.NewHandler(mm7cli, mmscli))
}
