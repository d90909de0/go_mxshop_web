package router

import (
	"github.com/gin-gonic/gin"
	"go_mxshop_web/user-web/api"
)

func InitUserRouter(router gin.RouterGroup) {
	userRouter := router.Group("user")
	userRouter.GET("list", api.GetUserList)
}
