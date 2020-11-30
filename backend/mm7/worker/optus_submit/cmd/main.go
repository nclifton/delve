package main

import (
	"fmt"
	"log"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"

	"github.com/burstsms/mtmo-tp/backend/mm7/worker"

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
	RPCHost               string `envconfig:"RPC_HOST"`
	RPCPort               int    `envconfig:"RPC_PORT"`
	OptusURL              string `envconfig:"OPTUS_URL"`
	OptusUser             string `envconfig:"OPTUS_USER"`
	OptusPass             string `envconfig:"OPTUS_PASSWORD"`
	TemplatePath          string `envconfig:"TEMPLATE_PATH"`
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
		RouteKey:      worker.QueueNameSubmitOptus,
		QueueName:     worker.QueueNameSubmitOptus,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)

	cli := client.NewClient(env.RPCHost, env.RPCPort)

	optusClient, err := tcl.NewService(env.OptusURL, env.OptusUser, env.OptusPass)
	if err != nil {
		log.Fatalf("failed to initialise optus client: %s reason: %s\n", Name, err)
	}

	soaptmpl := template.Must(template.ParseFiles(fmt.Sprintf(`%s/optus_submit.soap.tmpl`, env.TemplatePath)))

	w.Run(opts, optusWorker.NewHandler(cli, optusClient, soaptmpl))
}
