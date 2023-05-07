package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	HttpXBDCallerURI         = "X_BD_CALLER_URI"
	HttpXBDCallerURIV2       = "HTTP_X_BD_CALLER_URI"
	HttpUrlPressureCallerKey = "_caller_uri"
	HttpUrlPressureMarkKey   = "_press_mark"
)

func GetPressureFlag(ctx *gin.Context) (callerURI string, pressMark int) {
	if ctx != nil && ctx.Request != nil {
		callerURI = ctx.GetHeader(HttpXBDCallerURI)
		if callerURI == "" {
			callerURI = ctx.GetHeader(HttpXBDCallerURIV2)
		}
	}

	pressMark = 0
	if callerURI != "" && strings.Contains(callerURI, "/qa/test") {
		pressMark = 1
	}

	return callerURI, pressMark
}
