package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func NewLogger(name string) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	baseLogsPath := "./log/go-realworld/"
	if err := os.MkdirAll(baseLogsPath, 0777); err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s%s.log", baseLogsPath, name)
	_, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	cfg.OutputPaths = []string{filePath, "stdout"}
	return cfg.Build()
}
