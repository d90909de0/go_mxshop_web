package main

import (
	"fmt"
	"go_mxshop_web/user-web/global"
	"go_mxshop_web/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	initialize.Loggers()
	initialize.InitConfig()
	engine := initialize.Routers()

	port := global.ServerConfig.Port
	zap.S().Infof("启动服务器，端口：%d", port)
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
