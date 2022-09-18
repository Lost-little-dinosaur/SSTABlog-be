package captchaHandle

import (
	err2 "SSTABlog-be/internal/dto/err"
	"SSTABlog-be/internal/middleware"
	"SSTABlog-be/pkg/utils/captcha"
	"github.com/gin-gonic/gin"
)

func HandleGetCaptcha(c *gin.Context) {
	cp, err := captcha.GenerateCaptcha()
	if err != nil {
		middleware.Fail(c, err2.InternalErr)
		return
	}
	middleware.Success(c, *cp)
}
