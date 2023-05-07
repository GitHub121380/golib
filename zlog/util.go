package zlog

import (
	"fmt"
	"github.com/GitHub121380/golib/utils/metadata"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

const (
	ContextKeySpanID    = "spanId"
	ContextKeyLogID     = "logID"
	ContextKeyRequestID = "requestId"
)

func GetRequestID(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	requestId, exist := ctx.Get(ContextKeyRequestID)
	if exist {
		return requestId.(string)
	}
	return ""
}

// web 请求 兼容odp生成logid方式
func GetLogID(ctx *gin.Context) string {
	if ctx != nil {
		if logID := ctx.GetString(ContextKeyLogID); logID != "" {
			return logID
		}
		if ctx.Request != nil {
			if logID := ctx.GetHeader("X_BD_LOGID"); strings.TrimSpace(logID) != "" {
				ctx.Set(ContextKeyLogID, logID)
				return logID
			}
			if logID := ctx.GetHeader("x_bd_logid"); strings.TrimSpace(logID) != "" {
				ctx.Set(ContextKeyLogID, logID)
				return logID
			}
		}
	}

	usec := uint64(time.Now().UnixNano())
	logID := strconv.FormatUint(usec&0x7FFFFFFF|0x80000000, 10)

	// 这里有map并发写不安全问题，业务在job使用的时候规范ctx传参可避免，暂时不做加锁处理
	if ctx != nil {
		ctx.Set(ContextKeyLogID, logID)
	}

	return logID
}

func CreateSpan(ctx *gin.Context) string {
	return ""
}

func SetSpanID(ctx *gin.Context) string {
	return ""
}

// 兼容北斗的spanID生成方式
func GetSpanID(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	parentSpanID := ctx.GetString("parentSpanID")
	childSpanID := ctx.GetInt("childSpanID")
	newSpan := fmt.Sprintf("%s.%d", parentSpanID, childSpanID)
	return newSpan
}

// 用户自定义Notice
func AddNotice(ctx *gin.Context, key string, val interface{}) {
	if meta, ok := metadata.CtxFromGinContext(ctx); ok {
		if n := metadata.Value(meta, metadata.Notice); n != nil {
			if _, ok = n.(map[string]interface{}); ok {
				notices := n.(map[string]interface{})
				notices[key] = val
			}
		}
	}
}

// 获得所有用户自定义的Notice
func GetCustomerKeyValue(ctx *gin.Context) map[string]interface{} {
	meta, ok := metadata.CtxFromGinContext(ctx)
	if !ok {
		return nil
	}

	n := metadata.Value(meta, metadata.Notice)
	if n == nil {
		return nil
	}
	if notices, ok := n.(map[string]interface{}); ok {
		return notices
	}

	return nil
}

// server.log 中打印出用户自定义Notice
func PrintNotice(ctx *gin.Context) {
	notices := GetCustomerKeyValue(ctx)

	var fields []interface{}
	for k, v := range notices {
		fields = append(fields, k, v)
	}
	sugaredLogger(ctx).With(fields...).Info("notice")
}
