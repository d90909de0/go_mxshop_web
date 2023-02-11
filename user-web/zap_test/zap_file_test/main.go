package main

import (
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"./production.log",
		"stderr",
		"stdout",
	}

	return config.Build()
}

func main() {
	logger, _ := NewLogger()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	url := "http://www.imooc.con"
	sugar.Infof("Failed to fetch URL: %s", url)
}
