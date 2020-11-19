package rabbit_test

import (
	"testing"

	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
)

type MockHandler struct{}

// TODO actually make this do something
func (mck *MockHandler) Handle(body []byte, headers map[string]interface{}) error {
	return nil
}

func (mck *MockHandler) OnFinalFailure(body []byte) error {
	return nil
}

func TestRunWorker(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Error(err)
	}

	mockCon, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Error(err)
	}

	worker := rabbit.NewWorker("test_worker", mockCon, nil)

	go func() {
		worker.Run(rabbit.ConsumeOptions{
			PrefetchCount: 1,
			QueueName:     "testing",
			Exchange:      "testing",
			ExchangeType:  "direct",
			RouteKey:      "testing",
		}, &MockHandler{})
	}()

}
