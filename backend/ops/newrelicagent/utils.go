package newrelicagent

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Options holds config options for New Relic
type Options struct {
	AppName                       string
	NewRelicLicense               string
	DistributedTracerEnabled      bool
	AcceptDistributedTraceHeaders bool
}

// CreateApp creates a new newrelic application.
func CreateApp(opts *Options) *newrelic.Application {
	// New Relic is on by default apart from Local Dev where it is disabled via Envrionment Variable
	if os.Getenv("DISABLE_NEW_RELIC") == "true" {
		return nil
	}

	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("newrelic: create newrelic app with exception: %v", rec)
		}
	}()

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(opts.AppName),
		newrelic.ConfigLicense(opts.NewRelicLicense),
		newrelic.ConfigInfoLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(opts.DistributedTracerEnabled),
	)

	if err != nil {
		log.Printf("newrelic: failed to create newrelic app: %v", err)
		return nil
	}

	// Wait for the application to connect.
	if err = app.WaitForConnection(5 * time.Second); nil != err {
		log.Println(err)
		return nil
	}

	return app
}

// CreateAMQPHeader creates a AMQP formatted header to pass a reference for the given transaction onto a queue.
// This can then be consumed by another service when reading the job from the queue.
func CreateAMQPHeader(txn *newrelic.Transaction) map[string]interface{} {
	hdrs := http.Header{}
	txn.InsertDistributedTraceHeaders(hdrs)

	return map[string]interface{}{newrelic.DistributedTraceNewRelicHeader: hdrs.Get(newrelic.DistributedTraceNewRelicHeader)}
}

// AcceptAMQPHeader gets the AMQP formatted header from a job coming off a RabbitMQ queue.
// The header is then accepted by the passed in transaction.
// This stitches the transaction to a previous transaction and shows as one trace in New Relic.
func AcceptAMQPHeader(txn *newrelic.Transaction, headers map[string]interface{}) {
	traceHeader, ok := headers[newrelic.DistributedTraceNewRelicHeader].(string)

	if ok {
		hdrs := http.Header{}
		hdrs.Set(newrelic.DistributedTraceNewRelicHeader, traceHeader)
		txn.AcceptDistributedTraceHeaders(newrelic.TransportAMQP, hdrs)
	} else {
		log.Println("newrelic: amqp trace header does not exist")
	}
}

// StartAMQPSegment starts a MessageProducerSegment for monitoring jobs being added to a RabbitMQ queue.
func StartAMQPSegment(txn *newrelic.Transaction, exchangeName string) newrelic.MessageProducerSegment {
	seg := newrelic.MessageProducerSegment{
		Library:         "RabbitMQ",
		DestinationType: newrelic.MessageExchange,
		DestinationName: exchangeName,
	}

	seg.StartTime = txn.StartSegmentNow()

	return seg
}
