package baseServiceRouter

import (
	"SSTABlog-be/internal/handle/baseHandle/captchaHandle"
	"github.com/gin-gonic/gin"
)

func InitBaseServiceRouter(e *gin.Engine) {
	baseGroup := e.Group("/base")
	{
		baseGroup.GET("/captcha", captchaHandle.HandleGetCaptcha) //返回图片验证码图像的base64编码以及对应的id
	}
}
