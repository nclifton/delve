package main

import (
	"log"

	accountRPC "github.com/burstsms/mtmo-tp/backend/account/rpc"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort     int    `envconfig:"RPC_PORT"`
	PostgresURL string `envconfig:"POSTGRES_URL"`
}

func main() {
	var env Env
	err := envconfig.Process("account", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	arpc, err := accountRPC.NewService(env.PostgresURL)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", accountRPC.Name, err)
	}

	server, err := rpc.NewServer(arpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", accountRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", accountRPC.Name, port)
	server.Listen()
}
