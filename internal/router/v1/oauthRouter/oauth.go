package oauthRouter

import (
	"SSTABlog-be/internal/router/v1/oauthRouter/github"
	"github.com/gin-gonic/gin"
)

func InitOauthRouter(e *gin.Engine) {
	oauthGroup := e.Group("/oauth")
	github.InitGithubRouter(oauthGroup)
}
