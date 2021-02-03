package main

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/builder"
)

func main() {
	ctx := context.Background()
	s := rpcbuilder.NewGRPCServerFromEnv(ctx, builder.NewBuilderFromEnv())
	err := s.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
