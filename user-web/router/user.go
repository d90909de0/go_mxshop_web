package router

import (
	"github.com/gin-gonic/gin"
	"go_mxshop_web/user-web/api"
	"go_mxshop_web/user-web/middlewares"
)

func InitUserRouter(router gin.RouterGroup) {
	userRouter := router.Group("user")
	userRouter.GET("list", api.GetUserList).Use(middlewares.JWTAuth())
	userRouter.POST("login", api.PassWordLogin)
	userRouter.POST("add", api.CreateUser)
}
