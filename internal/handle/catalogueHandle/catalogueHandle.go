package catalogueHandle

import (
	"SSTABlog-be/internal/controller/articles"
	"SSTABlog-be/internal/controller/catalogues"
	"SSTABlog-be/internal/controller/users"
	"SSTABlog-be/internal/dto/article"
	"SSTABlog-be/internal/dto/catalogue"
	serviceErr "SSTABlog-be/internal/dto/err"
	"SSTABlog-be/internal/middleware"
	"SSTABlog-be/internal/model/Mysql"
	"github.com/gin-gonic/gin"
	"time"
)

func HandleAddCatalogue(c *gin.Context) {
	//登录验证
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	//权限验证
	if !users.PermissionCheck(uid, "1") { //需要1或0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}

	var req catalogue.AddCatalogueRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	//判断长度限制
	if len(req.CatalogueName) > 90 { //这里其实只能判断到字符串的长度，不能判断到中文的长度，中文只能存储20个字符
		middleware.FailWithCode(c, 40219, "目录名称过长")
		return
	}

	//查看同级目录名是否已存在
	if catalogues.CheckCatalogueExistByName(req.CatalogueName, req.FatherID) {
		middleware.FailWithCode(c, 40217, "同级目录下目录不可同名")
		return
	}

	//验证fatherID是否存在
	if len(req.FatherID) != 0 && catalogues.CheckCatalogueExist(req.FatherID) == nil {
		middleware.FailWithCode(c, 4029, "父级目录不存在")
		return
	}

	var tempString string
	err, tempString = catalogues.AddCatalogue(&Mysql.Catalogue{
		CatalogueName: req.CatalogueName,
		Description:   req.Description,
		CreateBy:      uid,
		FatherID:      req.FatherID,
		LastModifier:  uid,
	})
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, catalogue.GetCatalogueResponse{
		CatalogueID:      tempString,
		CatalogueName:    req.CatalogueName,
		Description:      req.Description,
		CreateBy:         uid,
		LastModifier:     uid,
		CreateOrUpdateAt: time.Now(),
		FatherID:         req.FatherID,
	})
	return
}
func HandleGetAllCatalogueSon(c *gin.Context) { //todo 增加返回排序
	catalogueID := c.Query("catalogueID")
	if len(catalogueID) == 0 {
		err, sonArr := GetAllCatalogueSon(catalogueID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, catalogue.GetCatalogueSonResponse{
			RootCatalogueID: "",
			CatalogueName:   "",
			Description:     "这里是顶层目录，其下不能有文章，只能有子目录",
			CreateBy:        "迷失的蓝色小恐龙",
			LastModifier:    "迷失的蓝色小恐龙",
			SonArr:          sonArr,
			//ArticleArr: //顶层目录下不能有文章
		})
	} else {
		if catalogues.CheckCatalogueExist(catalogueID) == nil {
			middleware.FailWithCode(c, 40218, "目录不存在")
			return
		}
		err, sonArr := GetAllCatalogueSon(catalogueID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		var tempStruct *Mysql.Catalogue
		tempStruct, err = catalogues.GetCatalogue(catalogueID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		var tempArticleArr []article.GetArticleInfoResponse
		err, tempArticleArr = GetArticlesByCatalogueID(tempStruct.ID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, catalogue.GetCatalogueSonResponse{
			RootCatalogueID:  tempStruct.ID,
			CatalogueName:    tempStruct.CatalogueName,
			LastModifier:     tempStruct.LastModifier,
			Description:      tempStruct.Description,
			CreateBy:         tempStruct.CreateBy,
			CreateOrUpdateAt: tempStruct.UpdatedAt,
			SonArr:           sonArr,
			ArticleArr:       tempArticleArr, //当前目录下的文章
		})
		return
	}
}
func HandleGetCatalogueSon(c *gin.Context) { //非递归获取子目录
	catalogueID := c.Query("catalogueID")
	if len(catalogueID) == 0 {
		err, sonArr := GetCatalogueSon(catalogueID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, catalogue.GetCatalogueSonResponse{
			RootCatalogueID: "",
			CatalogueName:   "",
			Description:     "这里是顶层目录，其下不能有文章，只能有子目录",
			CreateBy:        "迷失的蓝色小恐龙",
			LastModifier:    "迷失的蓝色小恐龙",
			SonArr:          sonArr,
			//ArticleArr: //顶层目录下不能有文章
		})
	} else {
		if catalogues.CheckCatalogueExist(catalogueID) == nil {
			middleware.FailWithCode(c, 40218, "目录不存在")
			return
		}
		err, sonArr := GetCatalogueSon(catalogueID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		var tempStruct *Mysql.Catalogue
		tempStruct, err = catalogues.GetCatalogue(catalogueID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		var tempArticleArr []article.GetArticleInfoResponse
		err, tempArticleArr = GetArticlesByCatalogueID(tempStruct.ID)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, catalogue.GetCatalogueSonResponse{
			RootCatalogueID:  tempStruct.ID,
			CatalogueName:    tempStruct.CatalogueName,
			LastModifier:     tempStruct.LastModifier,
			Description:      tempStruct.Description,
			CreateBy:         tempStruct.CreateBy,
			CreateOrUpdateAt: tempStruct.UpdatedAt,
			SonArr:           sonArr,
			ArticleArr:       tempArticleArr,
		})
		return
	}
}

func GetCatalogueSon(catalogueID string) (error, []catalogue.Son) { //非递归获取子目录
	var err error
	var tempCatalogue []Mysql.Catalogue
	var returnSonArr []catalogue.Son
	var tempArticleArr []article.GetArticleInfoResponse
	tempCatalogue, err = catalogues.GetCatalogueSons(catalogueID)
	if err != nil {
		return err, nil
	}
	for _, v := range tempCatalogue {
		err, tempArticleArr = GetArticlesByCatalogueID(v.ID)
		if err != nil {
			return err, nil
		}
		returnSonArr = append(returnSonArr, catalogue.Son{
			CatalogueID:      v.ID,
			CatalogueName:    v.CatalogueName,
			Description:      v.Description,
			CreateBy:         v.CreateBy,
			LastModifier:     v.LastModifier,
			CreateOrUpdateAt: v.UpdatedAt,
			ArticleArr:       tempArticleArr,
		})
		tempArticleArr = nil
	}
	return nil, returnSonArr
}

func GetAllCatalogueSon(catalogueID string) (error, []catalogue.Son) { //递归获取子目录
	var err error
	var tempCatalogue []Mysql.Catalogue
	var returnSonArr, tempArr []catalogue.Son
	var tempArticleArr []article.GetArticleInfoResponse
	tempCatalogue, err = catalogues.GetCatalogueSons(catalogueID)
	if err != nil {
		return err, nil
	}
	for _, v := range tempCatalogue {
		err, tempArr = GetAllCatalogueSon(v.ID)
		if err != nil {
			return err, nil
		}
		err, tempArticleArr = GetArticlesByCatalogueID(v.ID) //获取当前目录下的文章
		if err != nil {
			return err, nil
		}
		returnSonArr = append(returnSonArr, catalogue.Son{
			CatalogueID:      v.ID,
			CatalogueName:    v.CatalogueName,
			Description:      v.Description,
			CreateBy:         v.CreateBy,
			LastModifier:     v.LastModifier,
			CreateOrUpdateAt: v.UpdatedAt,
			SonArr:           tempArr,
			ArticleArr:       tempArticleArr,
		})
		tempArticleArr = nil
	}
	return nil, returnSonArr
}
func GetArticlesByCatalogueID(catalogueID string) (error, []article.GetArticleInfoResponse) {
	var tempArticleArr []article.GetArticleInfoResponse
	articleArr, err := articles.GetArticlesByCatalogueID(catalogueID)
	if err != nil {
		return err, nil
	}
	for _, vv := range articleArr {
		tempArticleArr = append(tempArticleArr, article.GetArticleInfoResponse{
			ID:            vv.ID,
			Title:         vv.Title,
			Cover:         vv.Cover,
			CreateBy:      vv.CreateBy,
			LastModifier:  vv.LastModifier,
			CatalogueID:   vv.CatalogueID,
			Description:   vv.Description,
			CommentNumber: vv.CommentNumber,
			PraiseNumber:  vv.PraiseNumber,
		})
	}
	return err, tempArticleArr
}

func HandleGetCatalogue(c *gin.Context) {
	catalogueID := c.Query("catalogueID")
	if len(catalogueID) == 0 {
		middleware.FailWithCode(c, 40218, "目录不存在")
		return
	}
	tempStruct := catalogues.CheckCatalogueExist(catalogueID)
	if tempStruct == nil {
		middleware.FailWithCode(c, 40218, "目录不存在")
		return
	}
	middleware.Success(c, catalogue.GetCatalogueResponse{
		CatalogueID:      tempStruct.ID,
		CatalogueName:    tempStruct.CatalogueName,
		Description:      tempStruct.Description,
		CreateBy:         tempStruct.CreateBy,
		LastModifier:     tempStruct.LastModifier,
		CreateOrUpdateAt: tempStruct.UpdatedAt,
		FatherID:         tempStruct.FatherID,
	})
	return
}

func HandleUpdateCatalogueName(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	if !users.PermissionCheck(uid, "1") { //需要1或0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	var req catalogue.UpdateCatalogueNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if len(req.CatalogueID) == 0 || catalogues.CheckCatalogueExist(req.CatalogueID) == nil {
		middleware.FailWithCode(c, 40218, "目录不存在")
		return
	}
	if len(req.CatalogueNewName) > 50 {
		middleware.FailWithCode(c, 40219, "目录名称过长")
		return
	}
	tempStruct, err := catalogues.GetCatalogue(req.CatalogueID)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	if tempStruct.CatalogueName == req.CatalogueNewName {
		middleware.Success(c, nil)
		return
	}
	if catalogues.CheckCatalogueExistByName(req.CatalogueNewName, tempStruct.FatherID) {
		middleware.FailWithCode(c, 40217, "同级目录下目录不可同名")
		return
	}
	err = catalogues.RenameCatalogue(req.CatalogueID, req.CatalogueNewName, uid)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
	return
}

func HandleUpdateCatalogueDescription(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	if !users.PermissionCheck(uid, "1") { //需要1或0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	var req catalogue.UpdateCatalogueDescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if len(req.CatalogueID) == 0 || catalogues.CheckCatalogueExist(req.CatalogueID) == nil {
		middleware.FailWithCode(c, 40218, "目录不存在")
		return
	}
	if len(req.NewDescription) > 255 {
		middleware.FailWithCode(c, 40220, "描述过长")
		return
	}
	err := catalogues.UpdateCatalogueDescription(req.CatalogueID, req.NewDescription, uid)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
	return
}

func HandleUpdateCatalogueFather(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	if !users.PermissionCheck(uid, "1") { //需要1或0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}
	var req catalogue.UpdateCatalogueParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if len(req.CatalogueID) == 0 || catalogues.CheckCatalogueExist(req.CatalogueID) == nil {
		middleware.FailWithCode(c, 40218, "目录不存在")
		return
	}
	tempStruct, err := catalogues.GetCatalogue(req.CatalogueID)
	if tempStruct.FatherID == req.NewFatherID { //如果目录的新父目录不变，则直接返回
		middleware.Success(c, nil)
		return
	}
	if catalogues.CheckCatalogueExistByName(tempStruct.CatalogueName, req.NewFatherID) {
		middleware.FailWithCode(c, 40217, "同级目录下目录不可同名")
		return
	}
	err = catalogues.UpdateCatalogueFather(req.CatalogueID, req.NewFatherID, uid)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
	return
}

func HandleDeleteCatalogue(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	if !users.PermissionCheck(uid, "1") { //需要1或0级权限
		middleware.FailWithCode(c, 40216, "对不起，您没有权限")
		return
	}

	catalogueID := c.Param("catalogueID")
	if len(catalogueID) == 0 || catalogues.CheckCatalogueExist(catalogueID) == nil {
		middleware.FailWithCode(c, 40218, "目录不存在")
		return
	}

	err := catalogues.DeleteCatalogue(catalogueID, uid)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	//删除目录后，删除目录下的所有文件
	var tempArticleArr []article.GetArticleInfoResponse
	err, tempArticleArr = GetArticlesByCatalogueID(catalogueID)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	for _, v := range tempArticleArr {
		err = articles.DeleteArticle(v.ID, uid)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
	}
	middleware.Success(c, nil)
	return
}

func HandleSearchCataloogue(c *gin.Context) { //todo 分页、文章查询
	keyWord := c.Query("keyWord")
	myType := c.Query("type")
	if len(keyWord) == 0 {
		middleware.FailWithCode(c, 40221, "搜索关键词不能为空")
		return
	}
	if myType == "description" {
		returnCatalogues, err := catalogues.SearchCatalogueByDescription(keyWord)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, returnCatalogues)
		return
	} else { //默认搜索名称
		returnCatalogues, err := catalogues.SearchCatalogueByName(keyWord)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, returnCatalogues)
		return
	}
}
