package main

import (
	"log"

	accountRPC "github.com/burstsms/mtmo-tp/backend/account/rpc"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort int `envconfig:"RPC_PORT"`
}

func main() {
	var env Env
	err := envconfig.Process("account", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	server, err := rpc.NewServer(accountRPC.NewService(), port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", accountRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", accountRPC.Name, port)
	server.Listen()
}
