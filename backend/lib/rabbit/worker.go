package rabbit

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/opentracing/opentracing-go"
)

type MessageHandler interface {
	Handle([]byte, map[string]interface{}) error
	OnFinalFailure([]byte) error
}

// TODO remove NR code litering our app
// should be replaced with calls to our own metrics service
type Worker struct {
	name   string
	con    Conn
	tracer opentracing.Tracer
}

func NewWorker(name string, con Conn, nr *nr.Options) *Worker {
	worker := &Worker{
		name: name,
		con:  con,
	}

	return worker
}

func NewWorkerWithTracer(name string, con Conn, nr *nr.Options, tracer opentracing.Tracer) *Worker {
	worker := &Worker{
		name:   name,
		con:    con,
		tracer: tracer,
	}

	return worker
}

func (w *Worker) Run(opts ConsumeOptions, handler MessageHandler) {
	log.Printf("%s worker started and waiting for jobs", opts.QueueName)

	messages, done, err := Consume(w.con, opts)
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
				log.Printf("received signal: %v, stopping rabbit consume", s)
				done <- true
			}
		}
	}()

	for msg := range messages {
		headers := amqpHeadersCarrier(msg.Headers())

		var sp opentracing.Span
		if w.tracer != nil {
			spCtx, err := w.tracer.Extract(opentracing.TextMap, headers)
			if err != nil {
				log.Printf("error parsing tracer span from message (%s): %s", msg.MessageId(), err)
			}

			sp = w.tracer.StartSpan(
				fmt.Sprintf("AMQP Consume %s %s", opts.Exchange, opts.RouteKey),
				opentracing.FollowsFrom(spCtx),
			)
			sp.LogKV("Message", msg.Body())
			// TODO: inject context with span into handler
		}

		err = handler.Handle(msg.Body(), headers)
		if err != nil {
			switch err.(type) {
			case *ErrWorkerMessageParse:
				log.Printf("error parsing message (%s): %s", msg.MessageId(), err)

				err := msg.Reject(false)
				if err != nil {
					log.Printf("could not reject message (%s): %s", msg.MessageId(), err)
				}
			case *ErrRetryWorkerMessage:
				log.Printf("error processing message from queue so requeing %s", err)

				if len(opts.RetryScale) > 0 {
					retryOpts, err := GenerateRetry(GenerateRetryOptions{
						Exchange:     fmt.Sprintf("%s-retry", opts.QueueName),
						ExchangeType: "topic",
						Delivery:     msg,
						MaxRetries:   len(opts.RetryScale),
						RouteKey:     opts.QueueName,
					})
					if err != nil {
						log.Printf("could not retry message (%s): %s", msg.MessageId(), err)
						err := handler.OnFinalFailure(msg.Body())
						if err != nil {
							log.Printf("could not handle final failure for message (%s): %s", msg.MessageId(), err)
						}

						err = msg.Reject(false)
						if err != nil {
							log.Printf("could not reject message (%s): %s", msg.MessageId(), err)
						}
						continue
					}
					err = Publish(w.con, retryOpts, msg.Body())
					if err != nil {
						log.Printf("could not retry message (%s): %s", msg.MessageId(), err)
					}
				} else {
					err := handler.OnFinalFailure(msg.Body())
					if err != nil {
						log.Printf("could not handle final failure for message (%s): %s", msg.MessageId(), err)
					}

					err = msg.Reject(true)
					if err != nil {
						log.Printf("could not reject message (%s): %s", msg.MessageId(), err)
					}
				}
			default:
				log.Printf("error processing message from queue %s ", err)
				err := handler.OnFinalFailure(msg.Body())
				if err != nil {
					log.Printf("could not handle final failure for message (%s): %s", msg.MessageId(), err)
				}
				err = msg.Reject(false)
				if err != nil {
					log.Printf("could not reject message (%s): %s", msg.MessageId(), err)
				}
			}
		}

		if err := msg.Ack(false); err != nil {
			log.Printf("failed to ack message (%s): %s", msg.MessageId(), err)
		}

		log.Printf("worker: %s successfully processed the msg", w.name)

		if w.tracer != nil {
			sp.Finish()
		}
	}

}
