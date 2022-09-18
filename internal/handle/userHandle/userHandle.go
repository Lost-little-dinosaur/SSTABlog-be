package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/users"
	serviceErr "github.com/wujunyi792/crispy-waffle-be/internal/dto/err"
	"github.com/wujunyi792/crispy-waffle-be/internal/dto/user"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
	"github.com/wujunyi792/crispy-waffle-be/internal/redis"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/jwtTokenGen"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/tecentCMS"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/captcha"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/check"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/crypto"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/gen/cmscode"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/gen/xrandom"
	"strings"
	"time"
)

func HandleSendRegisterCode(c *gin.Context) {
	var req user.SendCode
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
		middleware.FailWithCode(c, 40204, "验证码错误")
		return
	}

	if users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40205, "手机号已经注册，可以直接登录")
		return
	}

	_, err := redis.GetRedis().Get(req.Phone + "_register")
	if err == nil {
		middleware.FailWithCode(c, 40203, "发送过于频繁，稍后再试")
		return
	}
	data := cmscode.GenValidateCode(6)
	err = redis.GetRedis().Set(req.Phone+"_register", data, 5*time.Minute)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	tecentCMS.SendCMS(req.Phone, []string{"", data}) //这里要对应模板传递参数
	middleware.Success(c, nil)
}

func HandleCheckPhoneExist(c *gin.Context) {
	var req user.CheckPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	exist := users.CheckPhoneExist(req.Phone)
	middleware.Success(c, user.CheckPhoneResponse{
		Phone: req.Phone,
		Exist: exist,
	})
}

func HandleRegister(c *gin.Context) {
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40201, "手机号已存在")
		return
	}
	var err error
	var code string
	err = check.PasswordStrengthCheck(6, 20, 3, req.Password) //密码强度检查
	if err != nil {
		middleware.FailWithCode(c, 40202, err.Error())
		return
	}

	code, err = redis.GetRedis().Get(req.Phone + "_register")
	if err != nil || code != req.Code {
		middleware.Fail(c, serviceErr.CodeErr)
		return
	}

	err = redis.GetRedis().RemoveKey(req.Phone+"_register", false)
	if err != nil {
		return
	}

	salt := xrandom.GetRandom(5, xrandom.RAND_ALL)
	entity := Mysql.User{
		NickName:  "SSTA_" + xrandom.GetRandom(10, xrandom.RAND_NUM),
		Sex:       -1,
		Phone:     req.Phone,
		Signature: "这位用户没有任何想法~",
		Password:  crypto.PasswordGen(req.Password, salt),
		Salt:      salt,
	}
	err = users.RegisterUser(&entity)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, entity.Phone)
}

func HandleGeneralLogin(c *gin.Context) {
	var req user.LoginGeneralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if len(req.Info) != 11 && !strings.Contains(req.Info, "admin") { //如果是admin登录，则不需要验证手机号
		middleware.Fail(c, serviceErr.LoginErr)
		return
	}

	entity := Mysql.User{Phone: req.Info}
	users.GetEntity(&entity)

	if entity.ID == "" {
		middleware.Fail(c, serviceErr.LoginErr)
		return
	}

	if !crypto.PasswordCompare(req.Password, entity.Password, entity.Salt) {
		middleware.Fail(c, serviceErr.LoginErr)
		return
	}
	token, err := jwtTokenGen.GenToken(jwtTokenGen.Info{UID: entity.ID})
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, user.LoginResponse{Token: token})
	users.SetLoginLog(entity.ID, token)
}

func HandleSendPasswordResetCode(c *gin.Context) { //发送重置密码的验证码
	var req user.SendCode
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
		middleware.FailWithCode(c, 40204, "验证码错误")
		return
	}

	_, err := redis.GetRedis().Get(req.Phone + "_passwordReset")
	if err == nil {
		middleware.FailWithCode(c, 40203, "发送过于频繁，稍后再试")
		return
	}
	data := cmscode.GenValidateCode(6)
	err = redis.GetRedis().Set(req.Phone+"_passwordReset", data, 5*time.Minute)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	tecentCMS.SendCMS(req.Phone, []string{"", data})
	middleware.Success(c, nil)
}

func HandleResetPassword(c *gin.Context) {
	var req user.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if req.Password != req.PasswordAgain {
		middleware.FailWithCode(c, 40206, "密码不一致")
		return
	}
	if !users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40207, "手机号不存在")
		return
	}
	err := check.PasswordStrengthCheck(6, 20, 3, req.Password) //密码强度检查
	if err != nil {
		middleware.FailWithCode(c, 40202, err.Error())
		return
	}

	code, err := redis.GetRedis().Get(req.Phone + "_passwordReset")
	if err != nil || code != req.Code {
		middleware.Fail(c, serviceErr.CodeErr)
		return
	}

	_ = redis.GetRedis().RemoveKey(req.Phone+"_passwordReset", false)

	salt := xrandom.GetRandom(5, xrandom.RAND_ALL)
	passwordHashed := crypto.PasswordGen(req.Password, salt)
	err = users.UpdatePasswordAndSalt(req.Phone, passwordHashed, salt)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, req.Phone)
}

func HandleSendChangePhoneCode(c *gin.Context) { //发送修改手机号的验证码
	var req user.SendCode
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
		middleware.FailWithCode(c, 40204, "验证码错误")
		return
	}

	if users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40211, "手机号已被绑定")
		return
	}

	_, err := redis.GetRedis().Get(req.Phone + "_rebind")
	if err == nil {
		middleware.FailWithCode(c, 40203, "发送过于频繁，稍后再试")
		return
	}
	data := cmscode.GenValidateCode(6)
	err = redis.GetRedis().Set(req.Phone+"_rebind", data, 5*time.Minute)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	tecentCMS.SendCMS(req.Phone, []string{"", data})
	middleware.Success(c, nil)
}

func HandleUpdatePhone(c *gin.Context) {
	var req user.UpdatePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40201, "手机号已存在")
		return
	}
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	code, err := redis.GetRedis().Get(req.Phone + "_rebind")
	if err != nil || code != req.Code {
		middleware.Fail(c, serviceErr.CodeErr)
		return
	}

	_ = redis.GetRedis().RemoveKey(req.Phone+"_rebind", false)

	err = users.UpdatePhone(uid, req.Phone)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateNickName(c *gin.Context) { //修改昵称
	var req user.UpdateNickNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if users.CheckUserNameExist(req.NickName) {
		middleware.FailWithCode(c, 40208, "用户名已存在")
		return
	}
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	err := users.UpdateUserName(uid, req.NickName)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateAvatar(c *gin.Context) { //修改头像
	var req user.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	err := users.UpdateAvatar(uid, req.AvatarUrl)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateSex(c *gin.Context) {
	var req user.UpdateSexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	err := users.UpdateSex(uid, req.Sex)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateStatus(c *gin.Context) {
	var req user.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	err := users.UpdateStatus(uid, req.Status)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateSignature(c *gin.Context) {
	var req user.UpdateSignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	err := users.UpdateSignature(uid, req.Signature)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateEmail(c *gin.Context) {
	var req user.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if !check.VerifyEmailFormat(req.Email) {
		middleware.FailWithCode(c, 40209, "邮箱格式错误")
		return
	}

	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	err := users.UpdateEmail(uid, req.Email)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateStudentID(c *gin.Context) {
	var req user.UpdateStudentIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)

	if !check.VerifyStudentID(req.StudentID) { //如果检查不通过，则返回错误
		middleware.FailWithCode(c, 40215, "学号格式错误")
		return
	}

	err := users.UpdateStudentID(uid, req.StudentID)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleGetUserInfo(c *gin.Context) { //获取用户信息
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	entity := Mysql.User{}
	entity.ID = cuid.(string)
	users.GetEntity(&entity)
	permissions := users.GetPermission(uid)
	if permissions == nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, user.GetUserInfoResponse{
		Avatar:     entity.Avatar,
		Email:      entity.Email,
		StudentID:  entity.StudentID,
		Signature:  entity.Signature,
		RealName:   entity.RealName,
		NickName:   entity.NickName,
		Phone:      entity.Phone,
		Sex:        entity.Sex,
		Permission: permissions[0], //权限暂时只能有一个

	})
}

func HandleDelAccount(c *gin.Context) {
	middleware.Success(c, "功能暂不支持，如有需求，请联系工作人员")
}

func HandleAddPermission(c *gin.Context) { //TODO 修复外键错误
	var req user.AddPermissionRequest
	var err error
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if err = users.PermissionAdd(uid, req.PermissionName); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	//middleware.Success(c, "Success")

}

func HandleGetPermission(c *gin.Context) {
	cuid, _ := c.Get("uid")
	if cuid == nil {
		middleware.FailWithCode(c, 40214, "请先登录")
		return
	}
	uid := cuid.(string)
	permissions := users.GetPermission(uid)
	if permissions == nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, permissions)
	return
}
