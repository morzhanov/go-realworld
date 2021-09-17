package tracing

import (
	"fmt"

	"go.uber.org/zap"
)

type jaegerLogger struct {
	internalLogger *zap.Logger
}

type JaegerLogger interface {
	Error(msg string)
	Debug(msg string)
	Infof(msg string, args ...interface{})
}

func (l *jaegerLogger) Error(msg string) {
	l.internalLogger.Error(msg)
}

func (l *jaegerLogger) Debug(msg string) {
	l.internalLogger.Debug(msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	l.internalLogger.Info(fmt.Sprintf(msg, args))
}

func NewJeagerLogger(l *zap.Logger) JaegerLogger {
	return &jaegerLogger{internalLogger: l}
}
