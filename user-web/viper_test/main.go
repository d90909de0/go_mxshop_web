package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type MySqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// 嵌套配置
type ServerConfig struct {
	Name      string      `mapstructure:"name"`
	MySqlInfo MySqlConfig `mapstructure:"mysql"`
}

func main() {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("user-web/viper_test/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 读取配置到结构体
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)

	// 读取单个配置项
	fmt.Println(v.Get("name"))

	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(fmt.Sprintf("name:%s has changed", in.Name))
		fmt.Println(v.Get("name"))
	})

	time.Sleep(time.Second * 100)
}

// 读取环境变量【环境隔离】
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool("")
}
