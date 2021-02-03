package main

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/builder"
)

func main() {
	ctx := context.Background()
	s := rpcbuilder.NewGRPCServerFromEnv(ctx, builder.NewServiceFromEnv())

	err := s.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
