package server

import (
	"SSTABlog-be/config"
	_ "SSTABlog-be/internal/corn"
	"SSTABlog-be/internal/logger"
	"SSTABlog-be/internal/middleware"
	"SSTABlog-be/internal/redis"
	v1 "SSTABlog-be/internal/router/v1"
	"github.com/gin-gonic/gin"
)

var E *gin.Engine

func init() {
	logger.Info.Println("start init gin")
	gin.SetMode(config.GetConfig().MODE)
	E = gin.New()
	E.Use(middleware.GinRequestLog, gin.Recovery(), middleware.Cors(E))
}

func Run() {
	redis.GetRedis()
	v1.MainRouter(E)
	if err := E.Run("0.0.0.0:" + config.GetConfig().PORT); err != nil {
		logger.Error.Fatalln(err)
	}
}
