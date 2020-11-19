package rabbit

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/NeowayLabs/wabbit"
	"github.com/streadway/amqp"
)

// RetryScale is used as the default retry scale for retrying jobs
var RetryScale = []time.Duration{time.Minute, 2 * time.Minute, 5 * time.Minute, 10 * time.Minute, 30 * time.Minute}

// DeclareRetryQueues sets the resource up with the bindigns for a retry queue to a target DLX
func DeclareRetryQueues(con Conn, ResourceName string, RetryExchange string, RetryKey string, RetryScale []time.Duration) error {
	ch, err := con.Channel()
	if err != nil {
		return err
	}

	sourceExchange := fmt.Sprintf("%s-retry", ResourceName)

	if err = ch.ExchangeDeclare(
		sourceExchange,
		"topic",
		wabbit.Option{
			"durable":  true,
			"delete":   false,
			"internal": false,
			"noWait":   false,
		},
	); err != nil {
		return err
	}

	for i, time := range RetryScale {
		queue, err := ch.QueueDeclare(
			fmt.Sprintf("%s-retry%d", ResourceName, i+1),
			wabbit.Option{
				"durable":   true,
				"delete":    false,
				"exclusive": false,
				"noWait":    false,
				"args": amqp.Table{
					"x-dead-letter-exchange":    RetryExchange,
					"x-dead-letter-routing-key": RetryKey,
					"x-message-ttl":             time.Milliseconds(),
				},
			},
		)

		if err != nil {
			return err
		}

		bindKey := fmt.Sprintf("%s-retry%d", ResourceName, i+1)

		err = ch.QueueBind(queue.Name(), bindKey, sourceExchange, wabbit.Option{"noWait": false})

		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateRetryOptions specify headers and properties for a message before publishing.
type GenerateRetryOptions struct {
	// RouteKey specifies the routing key when publishing to exchange.
	RouteKey string

	// Exchange is the name of the exchange to publish.
	Exchange     string
	ExchangeType string

	Delivery Delivery

	MaxRetries int
}

var matchRetries = regexp.MustCompile(`-retry(\d+$)`)

// PublishRetry send a retry job with exponential backoff
func GenerateRetry(options GenerateRetryOptions) (PublishOptions, error) {
	if len(options.Delivery.Body()) <= 0 {
		return PublishOptions{}, fmt.Errorf("Cant retry a delivery wit an empty body: %+v", options.Delivery)
	}
	headers := options.Delivery.Headers()

	retryCount := 0

	// Use length of the x-death array to determine how many times job has been retried
	// See https://www.rabbitmq.com/dlx.html
	if headers["x-death"] != nil {
		death := headers["x-death"].([]interface{})
		if len(death) > 0 {
			retries := matchRetries.FindStringSubmatch(death[0].(amqp.Table)["queue"].(string))
			if retries != nil {
				retryCount, _ = strconv.Atoi(retries[1])
			}
		}
	}

	if retryCount >= options.MaxRetries {
		return PublishOptions{}, fmt.Errorf("Message exceeded > %d attempts so not retrying", options.MaxRetries)
	}

	retryCount++

	bindKey := fmt.Sprintf("%s-retry%d", options.RouteKey, retryCount)

	opt := PublishOptions{
		RouteKey:       bindKey,
		Exchange:       options.Exchange,
		ExchangeType:   options.ExchangeType,
		DontEncodeJson: true,
	}

	return opt, nil
}
