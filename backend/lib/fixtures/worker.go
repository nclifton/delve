package fixtures

import (
	"time"
)

func (tfx *TestFixtures) StartWorker(run func()) {

	// use go routine to run the webhook worker
	go run()

	// TODO see if we can use a health check here instead of a fixed wait time
	time.Sleep(100 * time.Millisecond) // force a wait a bit for the worker to become ready

}
