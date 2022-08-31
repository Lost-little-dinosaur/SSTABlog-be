package recycleBinRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/handle/recycleBin"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
)

func InitRecycleBinRouter(e *gin.Engine) {
	rB := e.Group("/recycleBin")
	rB.Use(middleware.JwtVerify) //需要登录才能访问
	{
		rB.GET("/getAllDelete", recycleBin.HandleGetRecycleBin)
		rB.DELETE("/restoreDelete", recycleBin.HandleRestoreDelete)
		rB.POST("/deleteForever", recycleBin.HandleDeleteForever) //todo 改请求方式

	}
}
