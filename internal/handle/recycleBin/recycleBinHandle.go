package recycleBin

import (
	"SSTABlog-be/internal/controller/articles"
	"SSTABlog-be/internal/controller/catalogues"
	"SSTABlog-be/internal/controller/users"
	serviceErr "SSTABlog-be/internal/dto/err"
	"SSTABlog-be/internal/dto/recycleBin"
	"SSTABlog-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func HandleGetRecycleBin(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	//权限验证
	if !users.PermissionCheck(uid, "3") { //需要3、2、1、0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	articleArr, err := articles.GetDeletedArticleInfo(uid)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	catalogueArr, err := catalogues.GetDeletedCatalogue(uid)
	middleware.Success(c, recycleBin.GetRecycleBinResponse{
		ArticleArr:   articleArr,
		CatalogueArr: catalogueArr,
	})
	return
}

func HandleRestoreDelete(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	//权限验证
	if !users.PermissionCheck(uid, "3") { //需要3、2、1、0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	var req recycleBin.RestoreDeleteRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if len(req.ID) == 0 {
		middleware.FailWithCode(c, 40226, "ID不能为空")
		return
	}
	if req.Type == "article" {
		err = articles.CheckIfArticleDeleted(uid, req.ID)
		if err != nil && err.Error() == "record not found" {
			middleware.FailWithCode(c, 40228, "回收站中找不到该文章")
			return

		} else if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		if err = articles.RestoreArticle(req.ID); err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, req.ID) //如果是文章的话，则返回ID以供预览
	} else if req.Type == "catalogue" {
		err = catalogues.CheckIfCatalogueDeleted(uid, req.ID)
		if err != nil && err.Error() == "record not found" {
			middleware.FailWithCode(c, 40228, "回收站中找不到该目录")
			return

		} else if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		if err = catalogues.RestoreCatalogue(req.ID); err != nil { //todo 判断原文件夹下是否存在同名目录
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, nil)
	} else {
		middleware.FailWithCode(c, 40227, "类型输入有误")
		return
	}
	return
}

func HandleDeleteForever(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	//权限验证
	if !users.PermissionCheck(uid, "3") { //需要3、2、1、0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	var req recycleBin.DeleteForeverRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if len(req.ID) == 0 {
		middleware.FailWithCode(c, 40226, "ID不能为空")
		return
	}
	if req.Type == "article" {
		err = articles.CheckIfArticleDeleted(uid, req.ID)
		if err != nil && err.Error() == "record not found" {
			middleware.FailWithCode(c, 40228, "回收站中找不到该文章")
			return

		} else if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		if err = articles.DeleteArticleForever(req.ID); err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, nil)
	} else if req.Type == "catalogue" {
		err = catalogues.CheckIfCatalogueDeleted(uid, req.ID)
		if err != nil && err.Error() == "record not found" {
			middleware.FailWithCode(c, 40228, "回收站中找不到该目录")
			return

		} else if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		if err = catalogues.DeleteCatalogueForever(req.ID); err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, nil)
	} else {
		middleware.FailWithCode(c, 40227, "类型输入有误")
		return
	}
	return
}
