package main

import (
	"fmt"
	"log"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"

	"github.com/burstsms/mtmo-tp/backend/mm7/worker"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	tcl "github.com/burstsms/mtmo-tp/backend/lib/tecloo/client"
	mm7w "github.com/burstsms/mtmo-tp/backend/mm7/worker/fake_mm7_submit"
	"github.com/kelseyhightower/envconfig"
)

var Name = "mm7-worker-fake-submit"

type Env struct {
	RabbitURL             string `envconfig:"RABBIT_URL"`
	RabbitExchange        string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType    string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RabbitPrefetchedCount int    `envconfig:"RABBIT_PREFETCHED_COUNT"`
	RPCHost               string `envconfig:"RPC_HOST"`
	RPCPort               int    `envconfig:"RPC_PORT"`
	TeclooURL             string `envconfig:"TECLOO_URL"`
	TemplatePath          string `envconfig:"TEMPLATE_PATH"`
}

func main() {
	log.Printf("starting worker: %s", Name)

	var env Env
	err := envconfig.Process("mm7", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise rabbit: %s reason: %s\n", Name, err)
	}

	opts := rabbit.ConsumeOptions{
		PrefetchCount: env.RabbitPrefetchedCount,
		Exchange:      env.RabbitExchange,
		ExchangeType:  env.RabbitExchangeType,
		RouteKey:      worker.QueueNameSubmitFake,
		QueueName:     worker.QueueNameSubmitFake,
	}

	w := rabbit.NewWorker(Name, rabbitmq, nil)

	cli := client.NewClient(env.RPCHost, env.RPCPort)

	tecloo, err := tcl.NewService(env.TeclooURL)
	if err != nil {
		log.Fatalf("failed to initialise tecloo client: %s reason: %s\n", Name, err)
	}

	soaptmpl := template.Must(template.ParseFiles(fmt.Sprintf(`%s/tecloo_submit.soap.tmpl`, env.TemplatePath)))

	w.Run(opts, mm7w.NewHandler(cli, tecloo, soaptmpl))
}
