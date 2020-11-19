package rabbit

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NeowayLabs/wabbit"
)

type Handler interface {
	Handle([]byte) error
	OnFinalFailure([]byte) error
}

type Consumer struct {
	name     string
	con      Conn
	queue    string
	prefetch int
}

type ConsumerOptions struct {
	Name          string
	Connection    Conn
	Queue         string
	PrefetchCount int
}

type ConsumerError struct {
	Message string
	Retry   bool
}

func (e *ConsumerError) Error() string {
	return e.Message
}

func NewConsumer(opts ConsumerOptions) (*Consumer, error) {
	// setup exchange

	// setup queues

	consumer := &Consumer{
		name:     opts.Name,
		con:      opts.Connection,
		queue:    opts.Queue,
		prefetch: opts.PrefetchCount,
	}

	return consumer, nil
}

func (c *Consumer) Run(opts ConsumeOptions, handler Handler) {
	log.Printf("started and waiting for jobs")

	messages, done, err := Consume(c.con, opts)
	if err != nil {
		log.Fatalf("failed to consume from queue: %s", err)
	}

	// listen for termination signals so we can cleanly close consumer
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		for s := range sig {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Printf("received signal: %v, stopping wabbit consume", s)
				done <- true
			}
		}
	}()

	for msg := range messages {
		log.Printf("processing message:")

		err = handler.Handle(msg.Body())
		if err != nil {
			switch err.(type) {
			case *ConsumerError:
				log.Printf("error reading json message (%q): %s", string(msg.Body()), err)

				err := msg.Reject(false)
				if err != nil {
					log.Printf("could not reject message (%q): %s", string(msg.Body()), err)
				}

			default:
				log.Printf("error processing message from queue %s ", err)

				err := handler.OnFinalFailure(msg.Body())
				if err != nil {
					log.Printf("could not handle final failure for message (%q): %s", string(msg.Body()), err)
				}

				err = msg.Reject(false)
				if err != nil {
					log.Printf("could not reject message (%q): %s", string(msg.Body()), err)
				}
			}
		}

		if err := msg.Ack(false); err != nil {
			log.Printf("failed to ack message (%q): %s", string(msg.Body()), err)
		}

		log.Printf("worker: %s successfully processed the msg", c.name)
	}
}

type ConsumeOptions struct {
	PrefetchCount int
	QueueName     string
	Exchange      string
	ExchangeType  string
	RouteKey      string
	RetryScale    []time.Duration
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

func closeChannel(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("could not close RabbitMQ channel: %s", err.Error())
	}
}
