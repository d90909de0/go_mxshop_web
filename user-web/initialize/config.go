package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_mxshop_web/user-web/config"
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
	nacosConfig := config.NacosConfig{}
	if err := v.Unmarshal(&nacosConfig); err != nil {
		zap.S().Panicf("服务配置文件初始化解析异常：%s", err.Error())
	}

	fmt.Println("nacosConfig：", nacosConfig)

	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosConfig.Nacos.NamespaceId, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      nacosConfig.Nacos.Host,
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-srv",
		Group:  "dev"})

	if err != nil {
		panic(err)
	}

	fmt.Println(content)

	if err = json.Unmarshal([]byte(content), global.ServerConfig); err != nil {
		panic(err)
	}

	// 读取单个配置项
	//fmt.Println(v.Get("name"))

	// 监听配置文件变化
	//v.WatchConfig()
	//v.OnConfigChange(func(in fsnotify.Event) {
	//	zap.S().Infof("服务配置变更", in.Name)
	//	if err := v.ReadInConfig(); err != nil {
	//		zap.S().Errorf("服务配置文件变更读取异常", err.Error())
	//	}
	//
	//	if err := v.Unmarshal(global.ServerConfig); err != nil {
	//		zap.S().Errorf("服务配置文件变更解析异常", err.Error())
	//	}
	//
	//	zap.S().Infof("服务配置变更后信息：%v", global.ServerConfig)
	//})
}
