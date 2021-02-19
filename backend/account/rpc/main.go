package main

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/builder"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
)

func main() {
	ctx := context.Background()
	s := rpcbuilder.NewGRPCServerFromEnv(ctx, builder.NewBuilderFromEnv())
	err := s.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
