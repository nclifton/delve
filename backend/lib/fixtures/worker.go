package fixtures

import (
	"log"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
)

func (tfx *TestFixtures) StartWorker(name string, workerService workerbuilder.Service) {

	worker := workerbuilder.NewWorker(
		workerbuilder.Config{
			WorkerName:                  name,
			RabbitURL:                   tfx.Rabbit.ConnStr,
			TracerDisable:               true,
			RabbitIgnoreClosedQueueConn: true,
			NRName:                      "",
			NRLicense:                   "",
			NRTracing:                   false,
		},
		workerService,
	)

	// use go routine to run the webhook worker
	go func() {
		if err := worker.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// TODO see if we can use a health check here instead of a fixed wait time
	time.Sleep(100 * time.Millisecond) // force a wait a bit for the worker to become ready

}
