package zlog

import (
	"github.com/GitHub121380/golib/env"
	"github.com/GitHub121380/golib/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetLogger 获得一个新的logger 会把日志打印到 name.log 中，不建议业务使用
// deprecated
func GetLogger() (s *zap.SugaredLogger) {
	if ServerLogger == nil {
		ServerLogger = newLogger("server").WithOptions(zap.AddCallerSkip(1)).Sugar()
	}
	return ServerLogger
}

// 通用字段封装
func sugaredLogger(ctx *gin.Context) *zap.SugaredLogger {
	if ctx == nil {
		return ServerLogger
	}

	return ServerLogger.With(
		zap.String("logId", GetLogID(ctx)),
		zap.String("spanId", GetSpanID(ctx)),
		zap.String("requestId", GetRequestID(ctx)),
		zap.String("module", env.AppName),
		zap.String("localIp", env.LocalIP),
		zap.String("handler", utils.GetHandler(ctx)), // todo: 后面优化去掉
	)
}

// 提供给业务使用的server log 日志打印方法
func Debug(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Debug(args...)
}

func Debugf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Debugf(format, args...)
}

func Info(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Info(args...)
}

func Infof(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Infof(format, args...)
}

func Warn(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Warn(args...)
}

func Warnf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Warnf(format, args...)
}

func Error(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Error(args...)
}

func Errorf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Errorf(format, args...)
}

func Panic(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Panic(args...)
}

func Panicf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Panicf(format, args...)
}

func Fatal(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Fatal(args...)
}

func Fatalf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Fatalf(format, args...)
}
