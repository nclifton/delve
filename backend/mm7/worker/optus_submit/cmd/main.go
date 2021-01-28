package main

import (
	"fmt"
	"log"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"

	"github.com/burstsms/mtmo-tp/backend/mm7/worker"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	tcl "github.com/burstsms/mtmo-tp/backend/lib/optus/client"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	optusWorker "github.com/burstsms/mtmo-tp/backend/mm7/worker/optus_submit"
	"github.com/kelseyhightower/envconfig"
)

var Name = "optus-submit"

type Env struct {
	RabbitURL             string `envconfig:"RABBIT_URL"`
	RabbitExchange        string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType    string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RabbitPrefetchedCount int    `envconfig:"RABBIT_PREFETCHED_COUNT"`
	MM7RPCAddress         string `envconfig:"MM7_RPC_ADDRESS"`
	OptusURL              string `envconfig:"OPTUS_URL"`
	OptusUser             string `envconfig:"OPTUS_USER"`
	OptusPass             string `envconfig:"OPTUS_PASSWORD"`
	TemplatePath          string `envconfig:"TEMPLATE_PATH"`

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
		RouteKey:      worker.QueueNameSubmitOptus,
		QueueName:     worker.QueueNameSubmitOptus,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)

	cli := client.NewClient(env.MM7RPCAddress)

	optusClient, err := tcl.NewService(env.OptusURL, env.OptusUser, env.OptusPass)
	if err != nil {
		log.Fatalf("Failed to initialise optus client: %s reason: %s\n", Name, err)
	}

	soaptmpl := template.Must(template.ParseFiles(fmt.Sprintf(`%s/optus_submit.soap.tmpl`, env.TemplatePath)))

	log.Println("Service started")
	w.Run(opts, optusWorker.NewHandler(cli, optusClient, soaptmpl))
}
