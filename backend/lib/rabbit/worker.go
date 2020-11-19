package rabbit

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	agent "github.com/burstsms/mtmo-tp/backend/ops/newrelicagent"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type MessageHandler interface {
	Handle([]byte, map[string]interface{}) error
	OnFinalFailure([]byte) error
}

// TODO remove NR code litering our app
// should be replaced with calls to our own metrics service
type Worker struct {
	name  string
	con   Conn
	nrApp *newrelic.Application
}

func NewWorker(name string, con Conn, nr *agent.Options) *Worker {
	worker := &Worker{
		name:  name,
		con:   con,
		nrApp: agent.CreateApp(nr),
	}

	return worker
}

func (w *Worker) Run(opts ConsumeOptions, handler MessageHandler) {
	log.Printf("started and waiting for jobs")

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
		nrTxn := w.nrApp.StartTransaction(w.name)
		log.Printf("processing message:")

		// Connect this transaction to previous ones in the chain
		agent.AcceptAMQPHeader(nrTxn, msg.Headers())

		err = handler.Handle(msg.Body(), msg.Headers())
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

		nrTxn.End()
	}
}
