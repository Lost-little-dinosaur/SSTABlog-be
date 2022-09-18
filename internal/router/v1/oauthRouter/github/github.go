package github

import (
	"SSTABlog-be/internal/handle/oauth/github"
	"SSTABlog-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitGithubRouter(e *gin.RouterGroup) {
	githubRouter := e.Group("/github")
	{
		githubRouter.GET("/login", github.HandleLogin)
		githubRouter.GET("/bind", middleware.JwtVerify, github.HandleBindAccount)
		githubRouter.GET("/callback", github.HandleCallBack)
	}
}
