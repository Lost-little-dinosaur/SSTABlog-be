package articleRouter

import (
	"SSTABlog-be/internal/handle/articleHandle"
	"SSTABlog-be/internal/middleware"
	"github.com/gin-gonic/gin"
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
		article.DELETE("/delete", articleHandle.HandleDeleteArticle)
	}
}
