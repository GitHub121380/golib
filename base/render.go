package base

import (
	"encoding/json"
	"fmt"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

// default render
type DefaultRender struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func RenderJson(ctx *gin.Context, code int, msg string, data interface{}) {
	renderJson := DefaultRender{code, msg, data}
	ctx.JSON(http.StatusOK, renderJson)
	//ctx.Set("render", renderJson)
	return
}

func RenderJsonSucc(ctx *gin.Context, data interface{}) {
	renderJson := DefaultRender{0, "succ", data}
	ctx.JSON(http.StatusOK, renderJson)
	//ctx.Set("render", renderJson)
	return
}

func RenderJsonFail(ctx *gin.Context, err error) {
	var renderJson DefaultRender

	switch errors.Cause(err).(type) {
	case Error:
		renderJson.ErrNo = errors.Cause(err).(Error).ErrNo
		renderJson.ErrMsg = errors.Cause(err).(Error).ErrMsg
		renderJson.Data = gin.H{}
	default:
		renderJson.ErrNo = -1
		renderJson.ErrMsg = errors.Cause(err).Error()
		renderJson.Data = gin.H{}
	}

	ctx.JSON(http.StatusOK, renderJson)
	//ctx.Set("render", renderJson)

	// 打印错误栈
	StackLogger(ctx, err)
	return
}

func RenderJsonAbort(ctx *gin.Context, err error) {
	var renderJson DefaultRender

	switch errors.Cause(err).(type) {
	case Error:
		renderJson.ErrNo = errors.Cause(err).(Error).ErrNo
		renderJson.ErrMsg = errors.Cause(err).(Error).ErrMsg
		renderJson.Data = gin.H{}
	default:
		renderJson.ErrNo = -1
		renderJson.ErrMsg = errors.Cause(err).Error()
		renderJson.Data = gin.H{}
	}

	ctx.AbortWithStatusJSON(http.StatusOK, renderJson)
	//ctx.Set("render", renderJson)

	return
}

// 打印错误栈
func StackLogger(ctx *gin.Context, err error) {
	if !strings.Contains(fmt.Sprintf("%+v", err), "\n") {
		return
	}

	var info []byte
	if ctx != nil {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorstack", "requestId": zlog.GetLogID(ctx)})
	} else {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorstack"})
	}

	fmt.Printf("%s\n-------------------stack-start-------------------\n%+v\n-------------------stack-end-------------------\n", string(info), err)
}
