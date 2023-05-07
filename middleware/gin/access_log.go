package gin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/GitHub121380/golib/base"
	"github.com/GitHub121380/golib/env"
	"github.com/GitHub121380/golib/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
)

type LoggerConfig struct {
	// requestBody 打印长度
	PrintRequestLen int
	// responsebody 打印长度
	PrintResponseLen int
	// mcpack数据协议的uri，请求参数打印原始二进制
	McpackReqUris []string
	// 请求参数不打印
	IgnoreReqUris []string
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	s = strings.Replace(s, "\n", "", -1)
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		//idx := len(b)
		// gin render json 后后面会多一个换行符
		//if b[idx-1] == '\n' {
		//	b = b[:idx-1]
		//}
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

var (
	// 暂不需要，后续考虑看是否需要支持用户配置
	mcpackReqUris []string
	ignoreReqUris []string
)

const (
	DefaultPrintRequestLen  = 10240
	DefaultPrintResponseLen = 10240
)

var printRequestLen, printResponseLen = -1, -1

func SetRequestAndResponseMaxShowLen(reqLen, respLen int) {
	printRequestLen = reqLen
	printResponseLen = respLen
}

func GetRequestAndResponseMaxShowLen() (int, int) {
	if printRequestLen < 0 {
		printRequestLen = DefaultPrintResponseLen
	}

	if printResponseLen < 0 {
		printResponseLen = DefaultPrintRequestLen
	}

	return printRequestLen, printResponseLen
}

// access日志打印
func AccessLog() gin.HandlerFunc {
	printRequestLen, printResponseLen := GetRequestAndResponseMaxShowLen()
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 请求url
		path := c.Request.URL.Path
		// 请求报文
		var requestBody []byte
		if c.Request.Body != nil {
			var err error
			requestBody, err = c.GetRawData()
			if err != nil {
				zlog.Warnf(c, "get http request body error: %s", err.Error())
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		blw := new(bodyLogWriter)
		if printResponseLen <= 0 {
			blw = &bodyLogWriter{body: nil, ResponseWriter: c.Writer}
		} else {
			blw = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		}
		c.Writer = blw

		c.Set("handler", c.HandlerName())
		logID := zlog.GetLogID(c)
		spanID := zlog.SetSpanID(c)

		// 处理请求
		c.Next()

		response := ""
		if blw.body != nil {
			if len(blw.body.String()) <= printResponseLen {
				response = blw.body.String()
			} else {
				response = blw.body.String()[:printResponseLen]
			}
		}

		bodyStr := ""
		flag := false
		// macpack的请求，以二进制输出日志
		for _, val := range mcpackReqUris {
			if strings.Contains(path, val) {
				bodyStr = fmt.Sprintf("%v", requestBody)
				flag = true
				break
			}
		}
		if !flag {
			// 不打印RequestBody的请求
			for _, val := range ignoreReqUris {
				if strings.Contains(path, val) {
					bodyStr = ""
					flag = true
					break
				}
			}
		}
		if !flag {
			bodyStr = string(requestBody)
		}

		if c.Request.URL.RawQuery != "" {
			bodyStr += "&" + c.Request.URL.RawQuery
		}

		if len(bodyStr) > printRequestLen {
			bodyStr = bodyStr[:printRequestLen]
		}

		// 结束时间
		end := time.Now()

		// 用户自定义notice
		var customerFields []zap.Field
		for k, v := range zlog.GetCustomerKeyValue(c) {
			customerFields = append(customerFields, zap.Reflect(k, v))
		}

		// 固定notice
		commonFields := []zap.Field{
			zap.String("logId", logID),
			zap.String("spanId", spanID),
			zap.String("requestId", zlog.GetRequestID(c)),
			zap.String("localIp", env.LocalIP),
			zap.String("module", env.AppName),
			zap.String("cuid", getReqValueByKey(c, "cuid")),
			zap.String("device", getReqValueByKey(c, "device")),
			zap.String("channel", getReqValueByKey(c, "channel")),
			zap.String("os", getReqValueByKey(c, "os")),
			zap.String("vc", getReqValueByKey(c, "vc")),
			zap.String("vcname", getReqValueByKey(c, "vcname")),
			zap.String("userid", getReqValueByKey(c, "userid")),
			zap.String("uri", c.Request.RequestURI),
			zap.String("host", c.Request.Host),
			zap.String("method", c.Request.Method),
			zap.String("httpProto", c.Request.Proto),
			zap.String("handle", c.HandlerName()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.String("refer", c.Request.Referer()),
			zap.String("clientIp", c.ClientIP()),
			zap.String("cookie", getCookie(c)),
			zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
			zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
			zap.Float64("cost", utils.GetRequestCost(start, end)),
			zap.String("requestParam", bodyStr),
			zap.Int("responseStatus", c.Writer.Status()),
			zap.String("response", response),
		}
		zlog.GetAccessLogger().With(commonFields...).Info("notice", customerFields...)
	}
}

// 从request body中解析特定字段作为notice key打印
func getReqValueByKey(ctx *gin.Context, k string) string {
	if vs, exist := ctx.Request.Form[k]; exist && len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func getCookie(ctx *gin.Context) string {
	cStr := ""
	for _, c := range ctx.Request.Cookies() {
		cStr += fmt.Sprintf("%s=%s&", c.Name, c.Value)
	}
	return strings.TrimRight(cStr, "&")
}

// access 添加kv打印
func AddNotice(k string, v interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		zlog.AddNotice(c, k, v)
		c.Next()
	}
}

func LoggerBeforeRun(ctx *gin.Context) {
	customCtx := ctx.CustomContext
	zlog.ServerLogger.With(
		"logId", zlog.GetLogID(ctx),
		"requestId", zlog.GetRequestID(ctx),
		"handle", customCtx.HandlerName(),
		"type", customCtx.Type,
	).Info("start")
}

func LoggerAfterRun(ctx *gin.Context) {
	customCtx := ctx.CustomContext
	cost := utils.GetRequestCost(customCtx.StartTime, customCtx.EndTime)
	var err error
	if customCtx.Error != nil {
		err = errors.Cause(customCtx.Error)
		base.StackLogger(ctx, customCtx.Error)
	}

	// 用户自定义notice
	notices := zlog.GetCustomerKeyValue(ctx)

	var fields []interface{}
	for k, v := range notices {
		fields = append(fields, k, v)
	}

	fields = append(fields,
		"logId", zlog.GetLogID(ctx),
		"requestId", zlog.GetRequestID(ctx),
		"handle", customCtx.HandlerName(),
		"type", customCtx.Type,
		"cost", cost,
		"error", fmt.Sprintf("%+v", err),
	)

	zlog.GetLogger().With(fields...).Info("end")
}
