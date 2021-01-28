package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/postbuilder"
)

func main() {
	s := workerbuilder.NewWorkerFromEnv(postbuilder.NewBuilderFromEnv())

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
