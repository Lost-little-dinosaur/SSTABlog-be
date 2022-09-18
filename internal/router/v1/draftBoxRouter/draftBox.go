package draftBoxRouter

import (
	"SSTABlog-be/internal/handle/draftHandle"
	"SSTABlog-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitDraftBoxRouter(e *gin.Engine) {
	draftBox := e.Group("/draftBox")
	draftBox.Use(middleware.JwtVerify) //需要登录才能访问
	{
		draftBox.POST("/add", draftHandle.HandleAddDraft)
		draftBox.POST("/update", draftHandle.HandleUpdateDraft)
		draftBox.DELETE("/delete", draftHandle.HandleDeleteDraftForever)
		draftBox.GET("/get", draftHandle.HandleGetDraft)
		draftBox.GET("/getAllInfo", draftHandle.HandleGetAllArticleInfo)
		draftBox.GET("/search", draftHandle.HandleSearchDraft)
	}
}
