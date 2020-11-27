package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	tlrpc "github.com/burstsms/mtmo-tp/backend/track_link/rpc"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RPCPort     int    `envconfig:"RPC_PORT"`
	PostgresURL string `envconfig:"POSTGRES_URL"`
	TrackHost   string `envconfig:"TRACK_HOST"`
}

func main() {
	var env Env
	err := envconfig.Process("track_link", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	port := env.RPCPort

	srpc, err := tlrpc.NewService(env.PostgresURL, env.TrackHost)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", tlrpc.Name, err)
	}

	server, err := rpc.NewServer(srpc, port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", tlrpc.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", tlrpc.Name, port)
	server.Listen()
}
