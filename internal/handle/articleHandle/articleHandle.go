package articleHandle

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/articles"
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/users"
	"github.com/wujunyi792/crispy-waffle-be/internal/dto/article"
	serviceErr "github.com/wujunyi792/crispy-waffle-be/internal/dto/err"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
)

func HandleAddArticle(c *gin.Context) {
	//登录验证
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

	var req article.AddArticleRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	//判断限制
	if len(req.Title) > 90 || len(req.Description) > 255 {
		middleware.FailWithCode(c, 40222, "标题或描述过长")
		return
	}
	if len(req.CatalogueID) == 0 {
		middleware.FailWithCode(c, 40223, "请选择目录")
		return
	}
	//查看同级目录同名文章是否已存在
	if articles.CheckArticleExistByCatalogueIDAndTitle(req.CatalogueID, req.Title) {
		middleware.FailWithCode(c, 40224, "同级目录下已存在同名文章")
		return
	}

	//添加文章
	var catalogueID string
	if err, catalogueID = articles.AddArticle(&Mysql.Article{
		Title:         req.Title,
		Description:   req.Description,
		CatalogueID:   req.CatalogueID,
		CreateBy:      uid,
		LastModifier:  uid,
		PraiseNumber:  0,
		CommentNumber: 0,
		WatchTimes:    0,
		Cover:         req.Cover,
		Content:       req.Content,
	}); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, catalogueID)
}

func HandleGetArticle(c *gin.Context) {
	articleID := c.Query("articleID")
	if len(articleID) == 0 {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	returnArticle, err := articles.GetArticleByID(articleID)
	if err != nil && err.Error() == "record not found" {
		middleware.FailWithCode(c, 40225, "找不到该文章")
		return
	} else if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, returnArticle)
	return
}

func HandleGetArticleInfo(c *gin.Context) {
	articleID := c.Query("articleID")
	if len(articleID) == 0 {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	returnArticle, err := articles.GetArticleByID(articleID)
	if err != nil && err.Error() == "record not found" {
		middleware.FailWithCode(c, 40225, "找不到该文章")
		return
	} else if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, article.GetArticleInfoResponse{
		Title:         returnArticle.Title,
		Cover:         returnArticle.Cover,
		CreateBy:      returnArticle.CreateBy,
		LastModifier:  returnArticle.LastModifier,
		CatalogueID:   returnArticle.CatalogueID,
		Description:   returnArticle.Description,
		CommentNumber: returnArticle.CommentNumber,
		PraiseNumber:  returnArticle.PraiseNumber,
	})
	return
}

func HandleSearchArticle(c *gin.Context) {
	keyword := c.Query("keyword")
	myType := c.Query("type")
	if len(keyword) == 0 {
		middleware.FailWithCode(c, 40221, "搜索关键词不能为空")
		return
	}
	if myType == "description" {
		articleArr, err := articles.SearchArticlesDescription(keyword)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, articleArr)
		return
	} else if myType == "content" {
		articleArr, err := articles.SearchArticlesContent(keyword)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, articleArr)
		return
	} else { //默认搜索标题
		articleArr, err := articles.SearchArticlesTitle(keyword)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, articleArr)
		return
	}
}

func HandleUpdateArticle(c *gin.Context) {
	//登录验证
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
	var req article.UpdateArticleRequest
	var err error
	var tempString string
	var tempFlag bool
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if len(req.ArticleID) == 0 {
		middleware.FailWithCode(c, 40225, "找不到该文章")
		return
	}
	tempFlag, tempString = articles.CheckArticleExistByID(req.ArticleID)
	if tempString != uid && !users.PermissionCheck(uid, "2") { //修改的文章不是自己的，需要2级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	if !tempFlag {
		middleware.FailWithCode(c, 40225, "找不到该文章")
		return
	}
	if len(req.Title) > 90 || len(req.Description) > 255 {
		middleware.FailWithCode(c, 40222, "标题或描述过长")
		return
	}
	//查看同级目录同名文章是否已存在
	var catalogueID, oldTitle string
	catalogueID, oldTitle, err = articles.GetArticleCatalogueIDAndTitleByID(req.ArticleID)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	if oldTitle != req.Title && articles.CheckArticleExistByCatalogueIDAndTitle(catalogueID, req.Title) {
		middleware.FailWithCode(c, 40224, "同级目录下已存在同名文章")
		return
	}
	if err = articles.UpdateArticle(&Mysql.Article{
		Title:        req.Title,
		Cover:        req.Cover,
		Description:  req.Description,
		LastModifier: uid,
		Content:      req.Content,
	}, req.ArticleID); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, req.ArticleID)
	return
}

func HandleDeleteArticle(c *gin.Context) {
	//登录验证
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
	articleID := c.Query("articleID")
	tempFlag, tempString := articles.CheckArticleExistByID(articleID)
	if len(articleID) == 0 || !tempFlag {
		middleware.FailWithCode(c, 40225, "找不到该文章")
		return
	}
	if tempString != uid && !users.PermissionCheck(uid, "2") { //删除的文章不是自己的，需要2级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	if err := articles.DeleteArticle(articleID, uid); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
	return
}
