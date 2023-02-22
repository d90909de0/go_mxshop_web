package main

import (
	"fmt"
	"go_mxshop_web/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	initialize.Loggers()
	engine := initialize.Routers()

	port := 8021
	zap.S().Infof("启动服务器，端口：%d", port)
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
