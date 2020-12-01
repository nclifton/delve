package main

import (
	"log"

	accountRPC "github.com/burstsms/mtmo-tp/backend/account/rpc"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort     int    `envconfig:"RPC_PORT"`
	PostgresURL string `envconfig:"POSTGRES_URL"`
	RedisURL    string `envconfig:"REDIS_URL"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("account", &env)
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

	port := env.RPCPort

	arpc, err := accountRPC.NewService(env.PostgresURL, env.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", accountRPC.Name, err)
	}

	server, err := rpc.NewServer(arpc, port)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", accountRPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", accountRPC.Name, port)
	server.Listen()
}
