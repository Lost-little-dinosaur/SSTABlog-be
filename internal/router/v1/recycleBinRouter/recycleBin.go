package recycleBinRouter

import (
	"SSTABlog-be/internal/handle/recycleBin"
	"SSTABlog-be/internal/middleware"
	"github.com/gin-gonic/gin"
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
