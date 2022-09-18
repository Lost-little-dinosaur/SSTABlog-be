package middleware

import (
	"SSTABlog-be/internal/dto/common"
	"SSTABlog-be/internal/dto/err"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Success 响应成功 ErrorCode 为 0 表示成功
func Success(c *gin.Context, data interface{}, count ...int64) {
	resp := common.JsonResponse{}
	resp.Clear()
	resp.Data = data
	if len(count) > 0 {
		resp.Count = count[0]
	}
	c.JSON(http.StatusOK, resp)
}

// Fail 响应失败
func Fail(c *gin.Context, serviceError err.ServiceError) {
	resp := common.JsonResponse{}
	resp.Clear()
	resp.Code = serviceError.Code
	resp.Message = serviceError.Error()
	c.JSON(serviceError.Code/100, resp)
	c.Abort()
}

func FailWithCode(c *gin.Context, code int, msg string) {
	resp := common.JsonResponse{}
	resp.Clear()
	resp.Message = msg
	resp.Code = code
	c.JSON(code/100, resp)
	c.Abort()
}
