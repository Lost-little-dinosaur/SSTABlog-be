package userRouter

import (
	"github.com/gin-gonic/gin"
	user "github.com/wujunyi792/crispy-waffle-be/internal/handle/userHandle"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
)

func InitUserRouter(e *gin.Engine) {
	userGroup := e.Group("/user")
	{
		userGroup.POST("/register/code", user.HandleSendRegisterCode) //发送注册验证码
		userGroup.POST("/register", user.HandleRegister)              //注册
		userGroup.POST("/exist/phone", user.HandleCheckPhoneExist)    //检查手机号是否存在
		userGroup.POST("/login/general", user.HandleGeneralLogin)     //普通登录

		userGroup.Use(middleware.JwtVerify) //需要登录才能访问
		{
			userGroup.GET("/info", user.HandleGetUserInfo)

			userGroup.POST("/addPermission", user.HandleAddPermission)
			userGroup.GET("/getPermission", user.HandleGetPermission)

			userGroup.POST("/resetPhone/code", user.HandleSendChangePhoneCode) //发送更换手机号验证码
			userGroup.POST("/resetPwd/code", user.HandleSendPasswordResetCode) //发送重置密码验证码
			userGroup.POST("/pwd/reset", user.HandleResetPassword)             //重置密码
			userGroup.POST("/reset/phone", user.HandleUpdatePhone)             //更换手机号

			userGroup.POST("/update/avatar", user.HandleUpdateAvatar)
			userGroup.POST("/update/nickname", user.HandleUpdateNickName)
			userGroup.POST("/update/studentid", user.HandleUpdateStudentID)
			userGroup.POST("/update/sex", user.HandleUpdateSex)
			userGroup.POST("/update/signature", user.HandleUpdateSignature)
			userGroup.POST("/update/status", user.HandleUpdateStatus)
			userGroup.POST("/update/email", user.HandleUpdateEmail)

			userGroup.DELETE("/delete", user.HandleDelAccount)
		}
	}
}
