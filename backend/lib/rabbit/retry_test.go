package rabbit_test

import (
	"testing"

	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/streadway/amqp"
)

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
				Delivery:     server.NewDelivery(nil, []byte("{}"), 1, "1", rabbit.Option{}),
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
				Delivery: server.NewDelivery(nil, []byte("{}"), 1, "1", rabbit.Option{
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
				Delivery: server.NewDelivery(nil, []byte("{}"), 1, "1", rabbit.Option{
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
				Delivery: server.NewDelivery(nil, []byte("{}"), 1, "1", rabbit.Option{
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
