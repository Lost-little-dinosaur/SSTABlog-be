package fileHandle

import (
	"SSTABlog-be/internal/middleware"
	"SSTABlog-be/internal/service/oss"
	"github.com/gin-gonic/gin"
)

func HandleGetAliUploadToken(c *gin.Context) {
	data := oss.GetPolicyToken()
	middleware.Success(c, data)
}

// HandleAliUpLoad 通过业务服务器中转文件至OSS 表单提交 字段名upload
func HandleAliUpLoad(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload") //upload为表单中文件字段名
	if err != nil {
		middleware.FailWithCode(c, 20008, err.Error())
	} else {
		url := oss.UploadFileToOss(header.Filename, file)
		if url == "" {
			middleware.FailWithCode(c, 50006, err.Error())
		}
		middleware.Success(c, url)
	}
}
