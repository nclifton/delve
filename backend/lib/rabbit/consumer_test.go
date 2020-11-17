package rabbit_test

import (
	"testing"

	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
)

type MockHandler struct{}

// TODO actually make this do something
func (mock *MockHandler) Handle(body []byte) error {
	return nil
}

func (mock *MockHandler) OnFinalFailure(body []byte) error {
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

	consumer, err := rabbit.NewConsumer(rabbit.ConsumerOptions{
		Name:          "test_worker",
		Connection:    mockCon,
		Queue:         "testing",
		PrefetchCount: 1,
	})
	if err != nil {
		t.Error(err)
	}

	go func() {
		consumer.Run(&MockHandler{})
	}()

}
