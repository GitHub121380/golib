package zlog

import (
	"github.com/GitHub121380/golib/env"
	"github.com/GitHub121380/golib/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 性能更高的logger
func GetAccessLogger() (l *zap.Logger) {
	if AccessLogger == nil {
		AccessLogger = newLogger(LogNameAccess)
	}
	return AccessLogger
}

func GetZapLogger() (l *zap.Logger) {
	if ModuleLogger == nil {
		ModuleLogger = newLogger(LogNameModule).WithOptions(zap.AddCallerSkip(1))
	}
	return ModuleLogger
}

func zapLogger(ctx *gin.Context) *zap.Logger {
	m := GetZapLogger()
	//m = m.WithOptions(zap.AddCallerSkip(1))
	if ctx == nil {
		return m
	}
	return m.With(zap.String("logId", GetLogID(ctx)),
		zap.String("spanId", GetSpanID(ctx)),
		zap.String("requestId", GetRequestID(ctx)),
		zap.String("module", env.GetAppName()),
		zap.String("localIp", env.LocalIP),
		zap.String("handler", utils.GetHandler(ctx)), // todo: 后面优化去掉
	)
}

func DebugLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Debug(msg, fields...)
}
func InfoLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Info(msg, fields...)
}

func WarnLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Warn(msg, fields...)
}

func ErrorLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Error(msg, fields...)
}

func PanicLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Panic(msg, fields...)
}

func FatalLogger(ctx *gin.Context, msg string, fields ...zap.Field) {
	zapLogger(ctx).Fatal(msg, fields...)
}
