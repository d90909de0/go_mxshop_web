package config

type ServerConfig struct {
	Name    string        `mapstructure:"name"`
	Port    int           `mapstructure:"port"`
	UserSrv UserSrvConfig `mapstructure:"userSrv"`
	JwtInfo JwtConfig     `mapstructure:"jwt"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key"`
}

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
