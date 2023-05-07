package golib

import (
	"github.com/GitHub121380/golib/base"
	"github.com/GitHub121380/golib/env"
	gg "github.com/GitHub121380/golib/middleware/gin"
	"github.com/gin-gonic/gin"
)

func Bootstrap(router *gin.Engine) {
	// 环境判断 env GIN_MODE=release/debug
	gin.SetMode(env.RunMode)

	// Global middleware
	router.Use(gg.Metadata())
	router.Use(gg.AccessLog())
	router.Use(gin.Recovery())

	// 存活探针
	router.HEAD("/health", base.HealthProbe())
	router.GET("/health", base.HealthProbe())

	// 就绪探针
	router.GET("/ready", base.ReadyProbe())

	// 性能分析工具
	base.Register(router)
}
