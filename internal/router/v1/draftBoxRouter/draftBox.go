package draftBoxRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/handle/draftHandle"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
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
