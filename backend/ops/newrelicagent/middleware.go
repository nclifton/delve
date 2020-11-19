package newrelicagent

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type newrelicmiddleware struct {
	handler                       http.Handler
	app                           *newrelic.Application
	acceptDistributedTraceHeaders bool
}

// ServeHTTP starts a new new relic transaction and injects it into the request chain.
// If nrm.app is nil then ServeHTTP does nothing and continues the chain.
// If nrm.acceptDistributedTraceHeaders is true then any New Relic Trace headers are accepted
// and Traces are stitched together into one.
func (nrm *newrelicmiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if nrm.app != nil {
		txn := nrm.app.StartTransaction(r.Method + " " + r.RequestURI)
		defer txn.End()

		w = txn.SetWebResponse(w)
		txn.SetWebRequestHTTP(r)

		if nrm.acceptDistributedTraceHeaders {
			txn.AcceptDistributedTraceHeaders(newrelic.TransportHTTPS, r.Header)
		}

		// Add transaction as header for services further down the chain
		txn.InsertDistributedTraceHeaders(r.Header)

		// Add transaction to context for use in this service
		r = newrelic.RequestWithTransactionContext(r, txn)
	}

	// Always route onwards even if newrelic didn't initialise
	// as newrelic not working won't affect end users.
	nrm.handler.ServeHTTP(w, r)
}

// New creates a new instance of the newrelic middleware.
// A new instance of a new relic application is created using the passed in options.
// If a new application instance fails to create then the returned app is set to nil.
// ServeHTTP() does nothing and continues the chain if app is nil and therefore doesn't break if new relic can't initialise.
func New(opts *Options) func(http.Handler) http.Handler {
	app := CreateApp(opts)
	return func(h http.Handler) http.Handler {
		return &newrelicmiddleware{
			handler:                       h,
			app:                           app,
			acceptDistributedTraceHeaders: opts.AcceptDistributedTraceHeaders,
		}
	}
}
