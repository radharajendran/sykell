package logger

import (
	"go.uber.org/zap"
)

var sugared *zap.SugaredLogger

func Init() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	logger, _ := cfg.Build()
	sugared = logger.Sugar()
}

func Sugar() *zap.SugaredLogger {
	return sugared
}