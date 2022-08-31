package draftHandle

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/draftBoxs"
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/users"
	"github.com/wujunyi792/crispy-waffle-be/internal/dto/draftBox"
	serviceErr "github.com/wujunyi792/crispy-waffle-be/internal/dto/err"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
)

func HandleAddDraft(c *gin.Context) {
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

	var req draftBox.AddDraftRequest
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

	//添加草稿
	if err = draftBoxs.AddDraft(&Mysql.Draft{
		Title:       req.Title,
		Description: req.Description,
		CreateBy:    uid,
		Cover:       req.Cover,
		Content:     req.Content,
	}); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil) //草稿不需要提供预览功能
}
func HandleGetDraft(c *gin.Context) {
	draftID := c.Query("draftID")
	if len(draftID) == 0 {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	returnDraft, err := draftBoxs.GetDraftByID(draftID)
	if err != nil && err.Error() == "record not found" {
		middleware.FailWithCode(c, 40225, "找不到该草稿")
		return
	} else if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, returnDraft)
	return
}

func HandleGetAllArticleInfo(c *gin.Context) { //todo 分页
	//登录验证
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	tempDraft, err := draftBoxs.GetAllArticleInfo(uid)
	if err != nil && err.Error() == "record not found" {
		middleware.FailWithCode(c, 40230, "草稿箱为空")
		return
	} else if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	returnDraft := make([]draftBox.GetAllDraftResponse, len(tempDraft))
	for i, v := range tempDraft {
		returnDraft[i] = draftBox.GetAllDraftResponse{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
		}
	}
	middleware.Success(c, returnDraft)
	return
}

func HandleUpdateDraft(c *gin.Context) {
	//登录验证
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	var req draftBox.UpdateDraftRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	var tempFlag bool
	var tempID string
	tempFlag, tempID = draftBoxs.CheckDraftExistByID(req.ID)
	if req.ID == "" || tempFlag == false || tempID != uid {
		middleware.FailWithCode(c, 40229, "找不到该草稿")
		return

	}
	//判断限制
	if len(req.Title) > 90 || len(req.Description) > 255 {
		middleware.FailWithCode(c, 40222, "标题或描述过长")
		return
	}
	//更新草稿
	if err = draftBoxs.UpdateDraft(&Mysql.Draft{
		Title:       req.Title,
		Description: req.Description,
		Cover:       req.Cover,
		Content:     req.Content,
	}, req.ID); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
	return
}

func HandleSearchDraft(c *gin.Context) {
	keyword := c.Query("keyword")
	myType := c.Query("type")
	if len(keyword) == 0 {
		middleware.FailWithCode(c, 40221, "搜索关键词不能为空")
		return
	}
	if myType == "description" {
		articleArr, err := draftBoxs.SearchDraftsDescription(keyword)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, articleArr)
		return
	} else if myType == "content" {
		articleArr, err := draftBoxs.SearchDraftsContent(keyword)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, articleArr)
		return
	} else { //默认搜索标题
		articleArr, err := draftBoxs.SearchDraftsTitle(keyword)
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, articleArr)
		return
	}
}

func HandleDeleteDraftForever(c *gin.Context) { //草稿要删除就是永久删除
	//登录验证
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	draftID := c.Query("draftID")
	if len(draftID) == 0 {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	tempFlag, tempUID := draftBoxs.CheckDraftExistByID(draftID)
	if tempFlag == false || tempUID != uid {
		middleware.FailWithCode(c, 40229, "找不到该草稿")
		return
	}
	if err := draftBoxs.DeleteDraftForever(draftID); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
	return
}
