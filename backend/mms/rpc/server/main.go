package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort     int    `envconfig:"RPC_PORT"`
	PostgresURL string `envconfig:"POSTGRES_URL"`
}

func main() {
	var env Env
	err := envconfig.Process("mms", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	arpc, err := mmsRPC.NewService(env.PostgresURL)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	server, err := rpc.NewServer(arpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", mmsRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", mmsRPC.Name, port)
	server.Listen()
}
