package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/builder"
)

func main() {
	s := rpcbuilder.NewGRPCServerFromEnv(builder.NewServiceFromEnv())

	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
