package rabbit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/streadway/amqp"
)

func TestDeclareRetryQueues(t *testing.T) {
	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	err := fakeServer.Start()
	if err != nil {
		t.Error(err)
	}

	mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Error(err)
	}

	testscale := []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}

	opts := rabbit.DLXOptions{
		ResourceName:   "test.thing",
		TargetExchange: "test.thing.exchange",
		TargetKey:      "test.thing.key",
		BackoffScale:   testscale,
		Type:           "retry",
	}
	err = rabbit.DeclareDLXQueues(mockConn, opts)
	if err != nil {
		t.Fatalf("Could not setup Retry Queues: %s", err)
	}

	channel, err := mockConn.Channel()
	if err != nil {
		t.Fatalf("Could not get mock channel: %s", err)
	}

	for i := 1; i <= len(testscale); i++ {
		queuename := fmt.Sprintf("test.thing-retry%d", i)
		_, err := channel.QueueInspect(queuename)
		if err != nil {
			t.Fatalf("Could not get expected retry queue: %s", queuename)
		}
	}

}

func TestGenerateRetry(t *testing.T) {

	retryTests := []struct {
		name          string
		expectedKey   string
		expectedError bool
		opts          rabbit.GenerateRetryOptions
	}{
		{
			name:        "initial retry",
			expectedKey: "testing-retry1",
			opts: rabbit.GenerateRetryOptions{
				RouteKey:     "testing",
				Exchange:     "testing-retry",
				ExchangeType: "topic",
				Delivery:     server.NewDelivery(nil, []byte("{}"), 1, "1", wabbit.Option{}),
				MaxRetries:   3,
			},
		},
		{
			name:        "second retry",
			expectedKey: "testing-retry2",
			opts: rabbit.GenerateRetryOptions{
				RouteKey:     "testing",
				Exchange:     "testing-retry",
				ExchangeType: "topic",
				Delivery: server.NewDelivery(nil, []byte("{}"), 1, "1", wabbit.Option{
					"x-death": []interface{}{
						amqp.Table{"queue": "testing-retry1"},
					},
				}),
				MaxRetries: 3,
			},
		},
		{
			name:        "third retry",
			expectedKey: "testing-retry3",
			opts: rabbit.GenerateRetryOptions{
				RouteKey:     "testing",
				Exchange:     "testing-retry",
				ExchangeType: "topic",
				Delivery: server.NewDelivery(nil, []byte("{}"), 1, "1", wabbit.Option{
					"x-death": []interface{}{
						amqp.Table{"queue": "testing-retry2"},
					},
				}),
				MaxRetries: 3,
			},
		},
		{
			name:          "forth retry",
			expectedKey:   "",
			expectedError: true,
			opts: rabbit.GenerateRetryOptions{
				RouteKey:     "testing",
				Exchange:     "testing-retry",
				ExchangeType: "topic",
				Delivery: server.NewDelivery(nil, []byte("{}"), 1, "1", wabbit.Option{
					"x-death": []interface{}{
						amqp.Table{"queue": "testing-retry3"},
					},
				}),
				MaxRetries: 3,
			},
		},
	}

	for _, test := range retryTests {
		t.Run(test.name, func(t *testing.T) {
			opts, err := rabbit.GenerateRetry(test.opts)
			if !test.expectedError && err != nil {
				t.Error(err)
			}
			if test.expectedError && err == nil {
				t.Fatalf("Expected an error from test: %+v not result; %+v", test, opts)
			}

			if opts.RouteKey != test.expectedKey {
				t.Fatalf("Did not get expected routekey (%s) for retry options: %+v", test.expectedKey, test.opts)
			}

		})
	}

}

func TestGenerateRequeue(t *testing.T) {

	retryTests := []struct {
		name          string
		expectedKey   string
		expectedError bool
		opts          rabbit.GenerateRequeueOptions
	}{
		{
			name:        "initial requeue",
			expectedKey: "testing-requeue1",
			opts: rabbit.GenerateRequeueOptions{
				RouteKey:     "testing",
				Exchange:     "testing-requeue",
				ExchangeType: "topic",
				Delivery:     server.NewDelivery(nil, []byte("{}"), 1, "1", wabbit.Option{}),
			},
		},
		{
			name:        "second requeue",
			expectedKey: "testing-requeue1",
			opts: rabbit.GenerateRequeueOptions{
				RouteKey:     "testing",
				Exchange:     "testing-requeue",
				ExchangeType: "topic",
				Delivery:     server.NewDelivery(nil, []byte("{}"), 1, "1", wabbit.Option{}),
			},
		},
	}

	for _, test := range retryTests {
		t.Run(test.name, func(t *testing.T) {
			opts, err := rabbit.GenerateRequeue(test.opts)
			if !test.expectedError && err != nil {
				t.Error(err)
			}
			if test.expectedError && err == nil {
				t.Fatalf("Expected an error from test: %+v not result; %+v", test, opts)
			}

			if opts.RouteKey != test.expectedKey {
				t.Fatalf("Did not get expected routekey (%s) for requeue options: %+v", test.expectedKey, test.opts)
			}

		})
	}

}
