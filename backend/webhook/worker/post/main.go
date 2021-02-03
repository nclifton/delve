package main

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/postbuilder"
)

func main() {
	ctx := context.Background()
	s := workerbuilder.NewWorkerFromEnv(ctx, postbuilder.NewBuilderFromEnv())

	if err := s.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
