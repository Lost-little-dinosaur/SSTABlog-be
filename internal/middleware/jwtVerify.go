package middleware

import (
	err2 "SSTABlog-be/internal/dto/err"
	"SSTABlog-be/internal/service/jwtTokenGen"
	"github.com/gin-gonic/gin"
)

func JwtVerify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		entry, err := jwtTokenGen.ParseToken(token)
		if err == nil {
			//c.Set("token", token)
			c.Set("uid", entry.Info.UID)
			c.Next()
			return
		} else {
			Fail(c, err2.JWTErr)
			return
		}
	}
	Fail(c, err2.VerifyErr)
	return
}
