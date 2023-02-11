package main

import (
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	url := "http://www.imooc.con"
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
	// sugar.Panic("Example")
	sugar.Error("Example")
	sugar.Debug("Example")

	// Logger
	logger.Info("failed to fetch RUL",
		zap.String("field", "value"),
		zap.Int("field", 1),
	)
}
