package tracing

import (
	"fmt"

	"go.uber.org/zap"
)

type JaegerLogger struct {
	internalLogger *zap.Logger
}

func (l *JaegerLogger) Error(msg string) {
	l.internalLogger.Error(msg)
}

func (l *JaegerLogger) Infof(msg string, args ...interface{}) {
	l.internalLogger.Info(fmt.Sprintf(msg, args))
}

func NewJeagerLogger(l *zap.Logger) *JaegerLogger {
	return &JaegerLogger{internalLogger: l}
}
