package rabbit_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"gotest.tools/assert"
)

func TestMessagePublish(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Error(err)
	}

	mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Error(err)
	}

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
		RouteKey:     "testing",
		Exchange:     "testing",
		ExchangeType: "direct",
	}

	err = rabbit.Publish(mockConn, testMessageOptions, "blah")
	if err != nil {
		t.Fatalf("Could not publish message: %s", err)
	}

}

func TestMessageConsume(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost.consume:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Fatal(err)
	}

	mockConn, err := amqptest.Dial("amqp://localhost.consume:5672/%2f")
	if err != nil {
		t.Error(err)
	}

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

	testJob := struct {
		Moose string `json:"moose"`
	}{Moose: "Marv"}

	testMessageOptions := rabbit.PublishOptions{
		RouteKey:     "testing",
		Exchange:     "testing",
		ExchangeType: "direct",
	}

	err = rabbit.Publish(mockConn, testMessageOptions, testJob)
	if err != nil {
		log.Printf("Could not publish message: %s", err)
	}
	expectedBody, err := json.Marshal(&testJob)
	if err != nil {
		t.Error(err)
	}

	message := <-deliveries
	assert.Equal(t, string(expectedBody), string(message.Body()))

}

func TestNonJsonMessagePublish(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost.consume.nonjson:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Fatal(err)
	}

	mockConn, err := amqptest.Dial("amqp://localhost.consume.nonjson:5672/%2f")
	if err != nil {
		t.Error(err)
	}

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

	testJob := []byte{12, 24, 45, 55}

	testMessageOptions := rabbit.PublishOptions{
		RouteKey:       "testing",
		Exchange:       "testing",
		ExchangeType:   "direct",
		DontEncodeJson: true,
	}

	err = rabbit.Publish(mockConn, testMessageOptions, testJob)
	if err != nil {
		log.Printf("Could not publish message: %s", err)
	}

	message := <-deliveries
	assert.Equal(t, string(testJob), string(message.Body()))

}
