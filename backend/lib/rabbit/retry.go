package rabbit

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/streadway/amqp"
)

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
	// See https://www.wabbitmq.com/dlx.html
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
		RouteKey: bindKey,
		Exchange: options.Exchange,
	}

	return opt, nil
}
