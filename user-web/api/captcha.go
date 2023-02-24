package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var captchaStore = base64Captcha.DefaultMemStore
var driver = base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)

func GetCaptcha(ctx *gin.Context) {
	id, s, err := base64Captcha.NewCaptcha(driver, captchaStore).Generate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "生成验证码错误"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId":   id,
		"captchaPath": s,
	})
}

func VerifyCaptcha(id, answer string) bool {
	return captchaStore.Verify(id, answer, true)
}
