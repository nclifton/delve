package rabbit

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
)

type Conn = wabbit.Conn
type Delivery = wabbit.Delivery
type Option = wabbit.Option

type PublishOptions struct {
	RouteKey       string
	Exchange       string
	ExchangeType   string
	Headers        map[string]interface{}
	Expiration     time.Duration
	Priority       uint8
	DontEncodeJson bool
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

	if options.DontEncodeJson {
		body = message.([]byte)
	} else {
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
