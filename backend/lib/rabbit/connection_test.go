package rabbit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
)

func TestMessagePublish(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Error(err)
	}
	defer fakeServer.Stop()

	mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Error(err)
	}
	defer mockConn.Close()

	_, _, err = rabbit.Consume(mockConn, rabbit.ConsumeOptions{
		PrefetchCount: 1,
		QueueName:     "testing",
		Exchange:      "testing",
		ExchangeType:  "direct",
		RouteKey:      "testing",
	})
	if err != nil {
		t.Fatalf("Could not setup consumer before publishing test message")
	}

	testMessageOptions := rabbit.PublishOptions{
		RouteKey: "testing",
		Exchange: "testing",
	}

	err = rabbit.Publish(mockConn, testMessageOptions, "blah")
	if err != nil {
		t.Fatalf("Could not publish message: %s", err)
	}
}

func TestMessageConsume(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer fakeServer.Stop()

	mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Error(err)
	}
	defer mockConn.Close()

	deliveries, _, err := rabbit.Consume(mockConn, rabbit.ConsumeOptions{
		PrefetchCount: 1,
		QueueName:     "testing",
		Exchange:      "testing",
		ExchangeType:  "direct",
		RouteKey:      "testing",
	})
	if err != nil {
		t.Fatalf("Could not open delivery channel: %s", err)
	}

	if len(deliveries) > 0 {
		t.Fatalf("Got unexpected deliveries in channel")
	}

}

func TestDeclareRetryQueues(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Error(err)
	}
	defer fakeServer.Stop()

	mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Error(err)
	}
	defer mockConn.Close()

	testscale := []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}

	err = rabbit.DeclareRetryQueues(mockConn, "test.thing", "test.thing.exchange", "test.thing.key", testscale)
	if err != nil {
		t.Fatalf("Could not setup Retry Queues: %s", err)
	}

	channel, err := mockConn.Channel()
	if err != nil {
		t.Fatalf("Could not get mock channel: %s", err)
	}
	defer channel.Close()

	for i := 1; i <= len(testscale); i++ {
		queuename := fmt.Sprintf("test.thing-retry%d", i)
		_, err := channel.QueueInspect(queuename)
		if err != nil {
			t.Fatalf("Could not get expected retry queue: %s", queuename)
		}
	}
}
