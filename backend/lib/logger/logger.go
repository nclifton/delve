package logger

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const spanId = "span-id"
const traceId = "trace-id"
const functionName = "function-name"

// TODO: replace fn args with runtime.Caller to log stack info?

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

type Fields map[string]interface{}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.TraceLevel)
	log.Debug("Logger created")

	return &StandardLogger{log}
}

// region Standard Log Interfaces
var (
	Info    = Event{1, "Info: %s"}
	Error   = Event{2, "Error: %s"}
	Warning = Event{3, "Warning: %s"}
	Debug   = Event{4, "Debug: %s"}
)

// WithFields
func (l *StandardLogger) Fields(ctx context.Context, fields Fields) *logrus.Entry {
	span, trace := getSpanAndTraceId(ctx)

	return l.WithFields(logrus.Fields{
		spanId:  span,
		traceId: trace,
	}).WithFields(logrus.Fields(fields))

}

// Info Logger
func (l *StandardLogger) Info(ctx context.Context, fn string, message string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: fn,
	}).Infof(Info.message, message)
}

// Infof Logger
func (l *StandardLogger) Infof(ctx context.Context, fn string, message string, args ...interface{}) {
	l.Info(ctx, fn, fmt.Sprintf(message, args...))
}

// Error Logger
func (l *StandardLogger) Error(ctx context.Context, fn string, message string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: fn,
	}).Errorf(Info.message, message)
}

// Errorf Logger
func (l *StandardLogger) Errorf(ctx context.Context, fn string, message string, args ...interface{}) {
	l.Error(ctx, fn, fmt.Sprintf(message, args...))
}

// Fatal Logger
func (l *StandardLogger) Fatal(ctx context.Context, fn string, message string) {
	l.Error(ctx, fn, message)
	os.Exit(1)
}

// Fatalf Logger
func (l *StandardLogger) Fatalf(ctx context.Context, fn string, message string, args ...interface{}) {
	l.Fatal(ctx, fn, fmt.Sprintf(message, args...))
}

// Warning Logger
func (l *StandardLogger) Warning(ctx context.Context, fn string, message string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: fn,
	}).Warningf(Info.message, message)
}

// Warningf Logger
func (l *StandardLogger) Warningf(ctx context.Context, fn string, message string, args ...interface{}) {
	l.Warning(ctx, fn, fmt.Sprintf(message, args...))
}

// Debug Logger
func (l *StandardLogger) Debug(ctx context.Context, fn string, message string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: fn,
	}).Debugf(Info.message, message)
}

// Debugf Logger
func (l *StandardLogger) Debugf(ctx context.Context, fn string, message string, args ...interface{}) {
	l.Debug(ctx, fn, fmt.Sprintf(message, args...))
}

// endregion

func getSpanAndTraceId(ctx context.Context) (spanId, traceId string) {
	sp := opentracing.SpanFromContext(ctx)

	if sp == nil {
		return "", ""
	}

	spanId = fmt.Sprint(sp)

	str := strings.Split(spanId, ":")
	if len(str) > 0 {
		traceId = str[0]
	}

	return
}
