package catalogueRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/handle/catalogueHandle"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
)

func InitCatalogueRouter(e *gin.Engine) {
	catalogues := e.Group("/catalogues")
	catalogues.GET("/listAll", catalogueHandle.HandleGetAllCatalogueSon) //递归获取所有目录
	catalogues.GET("/list", catalogueHandle.HandleGetCatalogueSon)       //非递归获取子目录
	catalogues.GET("/get", catalogueHandle.HandleGetCatalogue)
	catalogues.GET("/search", catalogueHandle.HandleSearchCataloogue)
	catalogues.Use(middleware.JwtVerify) //需要登录才能访问
	{
		catalogues.POST("/add", catalogueHandle.HandleAddCatalogue)
		catalogues.POST("/updateName", catalogueHandle.HandleUpdateCatalogueName)
		catalogues.POST("/updateDescription", catalogueHandle.HandleUpdateCatalogueDescription)
		catalogues.POST("/updateFather", catalogueHandle.HandleUpdateCatalogueFather)
		catalogues.POST("/deleteCatalogue", catalogueHandle.HandleDeleteCatalogue)
	}

}
