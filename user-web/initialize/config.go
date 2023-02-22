package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_mxshop_web/user-web/global"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("user-web/config/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf("服务配置文件初始化读取异常：%s", err.Error())
	}

	// 读取配置到结构体
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		zap.S().Panicf("服务配置文件初始化解析异常：%s", err.Error())
	}

	// 读取单个配置项
	fmt.Println(v.Get("name"))

	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("服务配置变更", in.Name)
		if err := v.ReadInConfig(); err != nil {
			zap.S().Errorf("服务配置文件变更读取异常", err.Error())
		}

		if err := v.Unmarshal(global.ServerConfig); err != nil {
			zap.S().Errorf("服务配置文件变更解析异常", err.Error())
		}

		zap.S().Infof("服务配置变更后信息：%v", global.ServerConfig)
	})
}
