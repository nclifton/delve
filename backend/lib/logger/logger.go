package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const spanId = "span-id"
const traceId = "trace-id"
const functionName = "function-name"

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
func (l *StandardLogger) Info(ctx context.Context, FuncName string, InfoMessage string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: FuncName,
	}).Infof(Info.message, InfoMessage)
}

// Error Logger
func (l *StandardLogger) Error(ctx context.Context, FuncName string, ErrorMessage string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: FuncName,
	}).Errorf(Info.message, ErrorMessage)
}

// Warning Logger
func (l *StandardLogger) Warning(ctx context.Context, FuncName string, WarningMessage string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: FuncName,
	}).Warningf(Info.message, WarningMessage)
}

// Debug Logger
func (l *StandardLogger) Debug(ctx context.Context, FuncName string, DebugMessage string) {
	span, trace := getSpanAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		spanId:       span,
		traceId:      trace,
		functionName: FuncName,
	}).Debugf(Info.message, DebugMessage)
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
