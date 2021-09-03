package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func NewLogger(name string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	filePath := fmt.Sprintf("/var/log/go-realworld/%v.log", name)
	cfg.OutputPaths = []string{filePath}
	return cfg.Build()
}
