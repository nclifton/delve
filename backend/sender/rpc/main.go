package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/builder"
)

func main() {
	s := rpcbuilder.NewGRPCServerFromEnv()
	b := builder.NewBuilderFromEnv()
	err := s.Start(b.Run)
	if err != nil {
		log.Fatal(err)
	}
}
