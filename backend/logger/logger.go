package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

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
func (l *StandardLogger) Fields(fields Fields) *logrus.Entry {
	return l.WithFields(logrus.Fields(fields))
}

// Info Logger
func (l *StandardLogger) Info(FuncName string, InfoMessage string) {
	l.Log(logrus.DebugLevel, "Starting InfoLog transaction now")
	e := l.WithContext(context.Background())
	e.Infof(Info.message, InfoMessage)
}

// Error Logger
func (l *StandardLogger) Error(FuncName string, ErrorMessage string) {
	l.Log(logrus.DebugLevel, "Starting ErrorLog transaction now")
	e := l.WithContext(context.Background())
	e.Errorf(Error.message, ErrorMessage)
}

// Warning Logger
func (l *StandardLogger) Warning(FuncName string, WarningMessage string) {
	l.Log(logrus.DebugLevel, "Starting WarningLog transaction now")
	e := l.WithContext(context.Background())
	e.Warningf(Warning.message, WarningMessage)
}

// Debug Logger
func (l *StandardLogger) Debug(FuncName string, DebugMessage string) {
	l.Log(logrus.DebugLevel, "Starting DebugLog transaction now")
	e := l.WithContext(context.Background())
	e.Debugf(Debug.message, DebugMessage)
}

// endregion
