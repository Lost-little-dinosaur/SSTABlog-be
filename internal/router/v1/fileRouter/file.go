package fileRouter

import (
	"SSTABlog-be/internal/handle/baseHandle/fileHandle"
	"github.com/gin-gonic/gin"
)

func InitFileRouter(e *gin.Engine) {
	file := e.Group("/fileHandle")
	{
		file.GET("/ali/token", fileHandle.HandleGetAliUploadToken)
		file.POST("/ali/upload", fileHandle.HandleAliUpLoad)
	}
}
