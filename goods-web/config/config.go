package config

type ServerConfig struct {
	Name     string         `mapstructure:"name" json:"name"`
	Port     int            `mapstructure:"port" json:"port"`
	GoodsSrv GoodsSrvConfig `mapstructure:"goodsSrv" json:"goodsSrv"`
	JwtInfo  JwtConfig      `mapstructure:"jwt" json:"jwt"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Nacos struct {
		Host        string `mapstructure:"host"`
		Port        uint64 `mapstructure:"port"`
		NamespaceId string `mapstructure:"namespaceId"`
		DataId      string `mapstructure:"dataId"`
		Group       string `mapstructure:"group"`
	} `mapstructure:"nacos"`
}
