package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
)

type WorkerConfig struct {
	ContainerName  string
	RabbitExchange string
	QueueName      string
}

func (tfx *TestFixtures) StartWorker(config WorkerConfig, workerService workerbuilder.Service) {

	port := tfx.workerPort()
	ctx := context.Background()
	worker := workerbuilder.NewWorker(ctx,
		workerbuilder.Config{
			ContainerName:               config.ContainerName,
			RabbitURL:                   tfx.Rabbit.ConnStr,
			TracerDisable:               true,
			RabbitIgnoreClosedQueueConn: true,
			NRName:                      "",
			NRLicense:                   "",
			NRTracing:                   false,
			QueueName:                   config.QueueName,
			RabbitExchange:              config.RabbitExchange,
			RabbitExchangeType:          "direct",
			RabbitPrefetchedCount:       1,
			HealthCheckHost:             tfx.env.HealthCheckHost,
			HealthCheckPort:             port,
			MaxGoRoutines:               200,
		},
		workerService,
	)
	tfx.WorkerHealthCheckURIs = append(tfx.WorkerHealthCheckURIs, fmt.Sprintf("http://%s:%s", tfx.env.HealthCheckHost, port))
	tfx.teardowns = append(tfx.teardowns, func() {
		worker.Stop(ctx)
	})

	// use go routine to run the webhook worker
	go func() {
		if err := worker.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// TODO see if we can use a health check here instead of a fixed wait time
	time.Sleep(100 * time.Millisecond) // force a wait a bit for the worker to become ready

}

func (tfx *TestFixtures) workerPort() string {
	tfx.workerPortIndex++
	if tfx.env.WorkerHealthCheckPorts[0] == "FREEPORT" || tfx.workerPortIndex > len(tfx.env.WorkerHealthCheckPorts) {
		return port(tfx.env.WorkerHealthCheckPorts[0])
	}
	return tfx.env.WorkerHealthCheckPorts[tfx.workerPortIndex]
}
