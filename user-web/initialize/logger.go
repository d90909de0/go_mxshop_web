package initialize

import "go.uber.org/zap"

func Loggers() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
