package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const requestId = "request-id"
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
	request, trace := getRequestAndTraceId(ctx)

	return l.WithFields(logrus.Fields{
		requestId: request,
		traceId:   trace,
	}).WithFields(logrus.Fields(fields))

}

// Info Logger
func (l *StandardLogger) Info(ctx context.Context, FuncName string, InfoMessage string) {
	request, trace := getRequestAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		requestId:    request,
		traceId:      trace,
		functionName: FuncName,
	}).Infof(Info.message, InfoMessage)
}

// Error Logger
func (l *StandardLogger) Error(ctx context.Context, FuncName string, ErrorMessage string) {
	request, trace := getRequestAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		requestId:    request,
		traceId:      trace,
		functionName: FuncName,
	}).Errorf(Info.message, ErrorMessage)
}

// Warning Logger
func (l *StandardLogger) Warning(ctx context.Context, FuncName string, WarningMessage string) {
	request, trace := getRequestAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		requestId:    request,
		traceId:      trace,
		functionName: FuncName,
	}).Warningf(Info.message, WarningMessage)
}

// Debug Logger
func (l *StandardLogger) Debug(ctx context.Context, FuncName string, DebugMessage string) {
	request, trace := getRequestAndTraceId(ctx)

	l.WithFields(logrus.Fields{
		requestId:    request,
		traceId:      trace,
		functionName: FuncName,
	}).Debugf(Info.message, DebugMessage)
}

// endregion

func getRequestAndTraceId(ctx context.Context) (requestId, traceId string) {
	sp := opentracing.SpanFromContext(ctx)

	requestId = fmt.Sprint(sp)

	str := strings.Split(requestId, ":")
	if len(str) > 0 {
		traceId = str[0]
	}

	return
}
