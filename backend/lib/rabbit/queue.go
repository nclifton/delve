package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/opentracing/opentracing-go"
	amqpT "github.com/streadway/amqp"
)

type Conn = wabbit.Conn
type Delivery = wabbit.Delivery
type Table = amqpT.Table

type PublishOptions struct {
	RouteKey       string
	Exchange       string
	ExchangeType   string
	Headers        map[string]interface{}
	Expiration     time.Duration
	Priority       uint8
	DontEncodeJson bool
	// TODO remove this dirt
	NrTxn  *newrelic.Transaction
	Tracer opentracing.Tracer
	Ctx    context.Context
}

type ConsumeOptions struct {
	PrefetchCount        int
	QueueName            string
	Exchange             string
	ExchangeType         string
	RouteKey             string
	RetryScale           []time.Duration
	AllowConnectionClose bool
}

// Publish will send each message specified by data to the exchange
// using the supplied routing key.
func Publish(con Conn, options PublishOptions, message interface{}) error {
	headers := amqpHeadersCarrier{}
	var err error
	var body []byte

	if options.DontEncodeJson {
		body = message.([]byte)
	} else {
		body, err = json.Marshal(&message)
		if err != nil {
			return err
		}
	}

	if options.Tracer != nil {
		parent := opentracing.SpanFromContext(options.Ctx)
		sp := options.Tracer.StartSpan(
			fmt.Sprintf("AMQP Publish %s %s", options.Exchange, options.RouteKey),
			opentracing.ChildOf(parent.Context()),
		)
		sp.LogKV("Message", string(body))
		sp.Finish()

		if err := options.Tracer.Inject(sp.Context(), opentracing.TextMap, headers); err != nil {
			return err
		}
	}

	ch, err := con.Channel()
	if err != nil {
		return err
	}

	// Declare the exchange
	// TODO move this out of publish, asking the server to confirm an exchange we know is there
	// thousands of times per second only serves to reduce performance
	// it also seems redundant even now given that worker consume declares it too
	if err = ch.ExchangeDeclare(
		options.Exchange,     // name of the exchange
		options.ExchangeType, // type
		wabbit.Option{
			"durable":  true,
			"delete":   false,
			"internal": false,
			"noWait":   false,
		},
	); err != nil {
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
		"headers":      Table(headers), // we DO need to convert to amqp.Table
		"deliveryMode": 2,              // 1 = Transient, 2 = Persistent
		"priority":     options.Priority,
	}

	if options.Expiration > 0 {
		// RabbitMQ expects Expiration specified in milliseconds
		params["expiration"] = strconv.FormatInt(options.Expiration.Milliseconds(), 10)
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

func Consume(con Conn, options ConsumeOptions) (chan Delivery, chan bool, error) {
	ch, err := con.Channel()
	if err != nil {
		return nil, nil, err
	}

	// Declare the exchange
	err = ch.ExchangeDeclare(
		options.Exchange,     // name of the exchange
		options.ExchangeType, // type
		wabbit.Option{
			"durable":  true,
			"delete":   false,
			"internal": false,
			"noWait":   false,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	// Declare Queue
	queue, err := ch.QueueDeclare(
		options.QueueName, // name of the queue
		wabbit.Option{
			"durable":   true,
			"delete":    false,
			"exclusive": false,
			"noWait":    false,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	// Bind the queue
	err = ch.QueueBind(
		queue.Name(),     // name of the queue
		options.RouteKey, // bindingKey
		options.Exchange, // sourceExchange
		wabbit.Option{
			"noWait": false,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	if len(options.RetryScale) > 0 {
		err = DeclareRetryQueues(con, options.QueueName, options.Exchange, options.RouteKey, options.RetryScale)
		if err != nil {
			return nil, nil, err
		}
	}

	err = ch.Qos(options.PrefetchCount, 0, false)
	if err != nil {
		defer closeChannel(ch)
		return nil, nil, err
	}

	c, err := ch.Consume(
		queue.Name(),
		"",
		wabbit.Option{
			"autoAck":   false,
			"exclusive": false,
			"noLocal":   false,
			"noWait":    false,
		},
	)
	if err != nil {
		defer closeChannel(ch)
		return nil, nil, err
	}

	dch := make(chan Delivery)
	done := make(chan bool)

	go func() {
		<-done
		log.Print("doing channel close on signal")
		closeChannel(ch)
	}()

	go func() {
		for d := range c {
			dch <- d
		}
		// Need to ensure we close this copied channel if c is closed (i.e. range ends)
		// so that if the server closes the connection/channel the caller
		// can handle this by knowing dch is closed
		close(dch)
	}()

	return dch, done, nil
}

/*
second argument optional and if is set true will allow closing of the connection without triggering a fatal
*/
func Connect(url string, flags ...bool) (Conn, error) {

	con, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	if len(flags) > 0 && flags[0] {
		return con, err
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
