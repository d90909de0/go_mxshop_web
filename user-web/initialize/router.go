package initialize

import (
	"github.com/gin-gonic/gin"
	"go_mxshop_web/user-web/router"
)

func Routers() *gin.Engine {
	engine := gin.Default()
	group := engine.RouterGroup
	router.InitUserRouter(group)
	return engine
}
