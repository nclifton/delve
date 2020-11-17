package logger

import (
	"bufio"
	"errors"
	"log"
	"strings"
	"time"

	"net"
	"net/http"
)

type logger struct {
	handler http.Handler
	opts    *Options
}

type loggerWriter struct {
	statusCode int
	http.ResponseWriter
}

func (lw *loggerWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *loggerWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := lw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("No rw hijacking from logger.")
	}
	return h.Hijack()
}

type Options struct {
	Verbose    bool
	UseXRealIp bool
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	// continue on but wrap the writer to capture written response code
	// Go sends 200 as default: http://golang.org/pkg/net/http/#ResponseWriter
	lw := &loggerWriter{200, w}
	l.handler.ServeHTTP(lw, r)

	if l.opts.Verbose {
		log.Printf("%s %s %s %d in %v", get_ip(r, l.opts.UseXRealIp), r.Method, r.RequestURI, lw.statusCode, time.Since(t))
	} else {
		log.Printf("%s in %v", r.RequestURI, time.Since(t))
	}
}

func New(opts *Options) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &logger{
			handler: h,
			opts:    opts}
	}
}

func get_ip(r *http.Request, use_xrealip bool) string {
	if use_xrealip && r.Header.Get("X-Real-Ip") != "" {
		return r.Header.Get("X-Real-Ip")
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
