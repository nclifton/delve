package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/run"
)

func main() {
	s := rpcbuilder.NewGRPCServerFromEnv()
	err := s.Start(run.Server)
	if err != nil {
		log.Fatal(err)
	}
}
