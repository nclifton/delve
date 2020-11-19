package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	receiver "github.com/burstsms/mtmo-tp/backend/tecloo_receiver"
	"github.com/kelseyhightower/envconfig"
)

var Name = "tecloo-receiver-api-http"

type Env struct {
	HTTPPort           int    `envconfig:"HTTP_PORT"`
	TemplatePath       string `envconfig:"TEMPLATE_PATH"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
}

type NREnv struct {
	Name    string `envconfig:"NAME"`
	License string `envconfig:"LICENSE"`
	Tracing bool   `envconfig:"TRACING"`
}

func main() {
	var env Env
	err := envconfig.Process("tecloo_receiver", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise rabbit: %s reason: %s\n", Name, err)
	}

	rabbitOpts := rabbit.PublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
	}

	var nrenv NREnv
	err = envconfig.Process("nr", &nrenv)
	if err != nil {
		log.Fatal("failed to read new relic env vars:", err)
	}

	port := strconv.Itoa(env.HTTPPort)

	newrelicM := nr.New(&nr.Options{
		AppName:                  nrenv.Name,
		NewRelicLicense:          nrenv.License,
		DistributedTracerEnabled: nrenv.Tracing,
	})

	opts := receiver.TeclooReceiverAPIOptions{
		NrApp:        newrelicM,
		TemplatePath: env.TemplatePath,
		Rabbit:       rabbitmq,
		RabbitOpts:   rabbitOpts,
	}

	server := receiver.NewTeclooReceiverAPI(&opts)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", Name, err)
	}

	log.Printf("%s service initialised and available on port %s", Name, port)
	log.Println("Tecloo Receiver API: listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))

}
