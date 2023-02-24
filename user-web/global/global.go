package global

import (
	ut "github.com/go-playground/universal-translator"
	"go_mxshop_web/user-web/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig *config.ServerConfig
	Trans        ut.Translator
)

func init() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:chengfei@20170917.com@tcp(116.62.167.224:3306)/go_mxshop_srvs_user?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	ServerConfig = &config.ServerConfig{}
}
