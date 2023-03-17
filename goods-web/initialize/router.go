package initialize

import (
	"github.com/gin-gonic/gin"
	"go_mxshop_web/goods-web/middlewares"
	"go_mxshop_web/goods-web/router"
)

func Routers() *gin.Engine {
	engine := gin.Default()
	engine.Use(middlewares.Cors())
	group := engine.RouterGroup
	router.InitGoodsRouter(group)
	return engine
}
