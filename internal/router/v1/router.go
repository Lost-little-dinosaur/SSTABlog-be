package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/config"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/articleRouter"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/baseServiceRouter"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/catalogueRouter"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/draftBoxRouter"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/fileRouter"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/oauthRouter"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/recycleBinRouter"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/userRouter"
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
			Test:       "test4",
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
