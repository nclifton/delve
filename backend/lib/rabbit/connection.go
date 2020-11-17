package rabbit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
)

type Delivery = wabbit.Delivery
type Option = wabbit.Option

type PublishOptions struct {
	RouteKey   string
	Exchange   string
	Headers    map[string]interface{}
	Expiration time.Duration
	Priority   uint8
}

// Publish will send each message specified by data to the exchange
// using the supplied routing key.
func Publish(con Conn, options PublishOptions, message interface{}) error {
	ch, err := con.Channel()
	if err != nil {
		return err
	}

	// Enabling confirm mode is a performance killer but very necessary as
	// this will ensure Rabbit reliably queues every message delivered.
	// Without confirm mode we can see performance of ~15000 msgs/s but with
	// a loss rate of about 0.2%. However in our business, a non-zero loss is unacceptable.
	// Confirm mode drops performance to about ~500 msgs/s at 0.0% loss.
	// The value passed in is mapped to "noWait", so noWait = false, being wait = true
	// meaning confirm mode is enabled by passing in false (dumb ass naming in this lib)
	err = ch.Confirm(false /*noWait*/)
	if err != nil {
		return err
	}
	defer closeChannel(ch)

	confirm := ch.NotifyPublish(make(chan wabbit.Confirmation, 1))

	params := wabbit.Option{
		"headers":      options.Headers, // we don't need to convert to amqp.Table (same type anyway)
		"deliveryMode": 2,               // 1 = Transient, 2 = Persistent
		"priority":     options.Priority,
	}

	if options.Expiration > 0 {
		// RabbitMQ expects Expiration specified in milliseconds
		params["expiration"] = strconv.FormatInt(options.Expiration.Milliseconds(), 10)
	}

	var body []byte

	// If we already have a string, assume its json and dont encode it twice
	switch message := message.(type) {
	case string:
		body = []byte(message)
	default:
		body, err = json.Marshal(&message)
		if err != nil {
			return err
		}

	}

	err = ch.Publish(options.Exchange, options.RouteKey, body, params)
	if err != nil {
		return fmt.Errorf("failed to publish message (%q): %s", string(body), err)
	}

	confirmation, ok := <-confirm
	if !ok {
		return errors.New("confirmation channel unexpectedly closed (likely amqp timeout)")
	}

	if !confirmation.Ack() {
		return errors.New("expecting ack after publishing to amqp")
	}

	return nil
}

type Conn = wabbit.Conn

func Connect(url string) (Conn, error) {
	con, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	// temporary solution to dropping amqp connections is to just make the process exit
	// let the service runners (k8s) handle restarting services
	conClosed := make(chan wabbit.Error)
	con.NotifyClose(conClosed)

	go func() {
		err := <-conClosed
		log.Fatalf("RabbitMQ Connection Closed: %s", err.Error())
	}()

	return con, nil
}

func closeChannel(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("could not close RabbitMQ channel: %s", err.Error())
	}
}

func DeclareExchange(con Conn, name, exchangeType string) error {
	ch, err := con.Channel()
	if err != nil {
		return err
	}
	defer closeChannel(ch)

	err = ch.ExchangeDeclare(
		name,
		exchangeType,
		wabbit.Option{
			"durable":  true,
			"delete":   false,
			"internal": false,
			"noWait":   false,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func DeclareQueue(con Conn, name, exchange, routeKey string, retryScales []time.Duration) error {
	ch, err := con.Channel()
	if err != nil {
		return err
	}

	// Declare Queue
	queue, err := ch.QueueDeclare(
		name,
		wabbit.Option{
			"durable":   true,
			"delete":    false,
			"exclusive": false,
			"noWait":    false,
		},
	)
	if err != nil {
		return err
	}

	// Bind the queue
	err = ch.QueueBind(
		queue.Name(),
		routeKey,
		exchange,
		wabbit.Option{
			"noWait": false,
		},
	)
	if err != nil {
		return err
	}

	if len(retryScales) > 0 {
		err = DeclareRetryQueues(con, queue.Name(), exchange, routeKey, retryScales)
		if err != nil {
			return err
		}
	}

	return nil
}

// the default scales for declaring retry queues
var RetryScales = []time.Duration{
	time.Minute,
	time.Minute * 2,
	time.Minute * 5,
	time.Minute * 10,
	time.Minute * 30,
}

// creates a new exchange appended with "-retry" and a set of queues bound to it
// it will create a queue for each scale and this queue will dead letter back to the
// given exchange using the routeKey
// https://www.wabbitmq.com/dlx.html may be useful reading
func DeclareRetryQueues(con Conn, queue, exchange, routeKey string, scales []time.Duration) error {
	ch, err := con.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	retryExchange := fmt.Sprintf("%s-retry", queue)

	err = ch.ExchangeDeclare(
		retryExchange,
		"topic",
		wabbit.Option{
			"durable":  true,
			"delete":   false,
			"internal": false,
			"noWait":   false,
		},
	)
	if err != nil {
		return err
	}

	// for each retry scale create a queue that sets a timeout for each message
	// these queues have no consumers and when messages hit their TTL they are
	// dead lettered back to the originating exchange for re-processing
	for i, ttl := range scales {
		queue, err := ch.QueueDeclare(
			fmt.Sprintf("%s%d", retryExchange, i+1),
			wabbit.Option{
				"durable":   true,
				"delete":    false,
				"exclusive": false,
				"noWait":    false,
				"args": wabbit.Option{
					"x-dead-letter-exchange":    exchange,
					"x-dead-letter-routing-key": routeKey,
					"x-message-ttl":             ttl.Milliseconds(),
				},
			},
		)
		if err != nil {
			return err
		}

		// now bind the retry queue created to the retry exchange
		// using the queue name as the routing key
		err = ch.QueueBind(queue.Name(), queue.Name(), retryExchange, wabbit.Option{"noWait": false})
		if err != nil {
			return err
		}
	}

	return nil
}
