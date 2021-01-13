package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/servicebuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/run"
)

func main() {
	s := servicebuilder.NewGRPCServerFromEnv()
	err := s.GRPCStart(run.Server)
	if err != nil {
		log.Fatal(err)
	}
}
