package articleRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/handle/articleHandle"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
)

func InitArticleRouter(e *gin.Engine) {
	article := e.Group("/article")
	article.GET("/get", articleHandle.HandleGetArticle)
	article.GET("/getInfo", articleHandle.HandleGetArticleInfo)
	article.GET("/search", articleHandle.HandleSearchArticle)
	article.Use(middleware.JwtVerify) //需要登录才能访问
	{
		article.POST("/update", articleHandle.HandleUpdateArticle)
		article.POST("/add", articleHandle.HandleAddArticle)
		article.POST("/delete", articleHandle.HandleDeleteArticle)
	}
}
