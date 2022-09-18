package v1

import (
	"SSTABlog-be/config"
	"SSTABlog-be/internal/middleware"
	"SSTABlog-be/internal/router/v1/articleRouter"
	"SSTABlog-be/internal/router/v1/baseServiceRouter"
	"SSTABlog-be/internal/router/v1/catalogueRouter"
	"SSTABlog-be/internal/router/v1/draftBoxRouter"
	"SSTABlog-be/internal/router/v1/fileRouter"
	"SSTABlog-be/internal/router/v1/oauthRouter"
	"SSTABlog-be/internal/router/v1/recycleBinRouter"
	"SSTABlog-be/internal/router/v1/userRouter"
	"fmt"
	"github.com/gin-gonic/gin"
)

func MainRouter(e *gin.Engine) {
	e.Any("", func(c *gin.Context) {
		data := struct {
			UA         string
			Host       string
			Method     string
			Proto      string
			RemoteAddr string
			Message    string
			Test       string
		}{
			UA:         c.Request.Header.Get("User-Agent"),
			Host:       c.Request.Host,
			Method:     c.Request.Method,
			Proto:      c.Request.Proto,
			RemoteAddr: c.Request.RemoteAddr,
			Message:    fmt.Sprintf("Welcome to %s, version %s.", config.GetConfig().ProgramName, config.GetConfig().VERSION),
			Test:       "test17",
		}
		middleware.Success(c, data)
	})
	baseServiceRouter.InitBaseServiceRouter(e)
	fileRouter.InitFileRouter(e)
	userRouter.InitUserRouter(e)
	oauthRouter.InitOauthRouter(e)
	catalogueRouter.InitCatalogueRouter(e)
	articleRouter.InitArticleRouter(e)
	recycleBinRouter.InitRecycleBinRouter(e)
	draftBoxRouter.InitDraftBoxRouter(e)
}
