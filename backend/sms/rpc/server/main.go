package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	smsRPC "github.com/burstsms/mtmo-tp/backend/sms/rpc"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort            int    `envconfig:"RPC_PORT"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
}

func main() {
	var env Env
	err := envconfig.Process("sms", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	rabbitOpts := smsRPC.RabbitPublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: "direct",
	}

	arpc, err := smsRPC.NewService(env.PostgresURL, rabbitmq, rabbitOpts)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	server, err := rpc.NewServer(arpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", smsRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", smsRPC.Name, port)
	server.Listen()
}
