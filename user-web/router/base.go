package router

import (
	"github.com/gin-gonic/gin"
	"go_mxshop_web/user-web/api"
)

func InitBaseRouter(router gin.RouterGroup) {
	baseRouter := router.Group("base")
	baseRouter.GET("captcha", api.GetCaptcha)
}
