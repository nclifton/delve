package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/postbuilder"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// TODO: use this code when dockerised webhook-post-worker-service is merged
	// s := workerbuilder.NewWorkerFromEnv(postbuilder.NewBuilderFromEnv())
	// ............

	// and remove this:
	config := struct {
		ContainerName string `envconfig:"CONTAINER_NAME"`
		RabbitURL     string `envconfig:"RABBIT_URL"`

		TracerDisable               bool `envconfig:"TRACER_DISABLE"`
		RabbitIgnoreClosedQueueConn bool `envconfig:"RABBIT_IGNORE_CLOSED_QUEUE_CONN"`

		NRName    string `envconfig:"NR_NAME"`
		NRLicense string `envconfig:"NR_LICENSE"`
		NRTracing bool   `envconfig:"NR_TRACING"`
	}{}
	if err := envconfig.Process("webhook", &config); err != nil {
		log.Fatal(err)
	}
	s := workerbuilder.NewWorker(config, postbuilder.NewBuilderFromEnv())
	// ..............

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
