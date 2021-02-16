package main

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker/mgage_submit/mgagesubmitbuilder"
)

func main() {
	ctx := context.Background()
	s := workerbuilder.NewWorkerFromEnv(ctx, mgagesubmitbuilder.NewBuilderFromEnv())

	if err := s.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
