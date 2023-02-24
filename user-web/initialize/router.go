package initialize

import (
	"github.com/gin-gonic/gin"
	"go_mxshop_web/user-web/middlewares"
	"go_mxshop_web/user-web/router"
)

func Routers() *gin.Engine {
	engine := gin.Default()
	engine.Use(middlewares.Cors())
	group := engine.RouterGroup
	router.InitUserRouter(group)
	router.InitBaseRouter(group)
	return engine
}
