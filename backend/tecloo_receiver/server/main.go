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

var Name = "tecloo-receiver"

type Env struct {
	HTTPPort           int    `envconfig:"HTTP_PORT"`
	TemplatePath       string `envconfig:"TEMPLATE_PATH"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`

	Name    string `envconfig:"NAME"`
	License string `envconfig:"LICENSE"`
	Tracing bool   `envconfig:"TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("tecloo_receiver", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	newrelicM := nr.New(&nr.Options{
		AppName:                  env.Name,
		NewRelicLicense:          env.License,
		DistributedTracerEnabled: env.Tracing,
	})

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise rabbit: %s reason: %s\n", Name, err)
	}

	rabbitOpts := rabbit.PublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
	}

	port := strconv.Itoa(env.HTTPPort)

	opts := receiver.TeclooReceiverAPIOptions{
		NrApp:        newrelicM,
		TemplatePath: env.TemplatePath,
		Rabbit:       rabbitmq,
		RabbitOpts:   rabbitOpts,
	}

	server := receiver.NewTeclooReceiverAPI(&opts)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", Name, err)
	}

	log.Printf("%s service initialised and available on port %s", Name, port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
