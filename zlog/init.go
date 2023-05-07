package zlog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/GitHub121380/golib/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 对用户暴露的log配置
type Rotate struct {
	Switch bool   `yaml:"switch"`
	Unit   string `yaml:"unit"`
	Count  int    `yaml:"count"`
}
type LogConfig struct {
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
	Rotate Rotate `yaml:"rotate"`
}

type loggerConfig struct {
	ZapLevel zapcore.Level

	// 以下变量仅对开发环境生效
	Stdout   bool
	Log2File bool
	Path     string

	RotateUnit   string
	RotateCount  int
	RotateSwitch bool
}

// 全局配置 仅限Init函数进行变更
var logConfig = loggerConfig{
	ZapLevel: zapcore.InfoLevel,

	Stdout:   false,
	Log2File: true,
	Path:     "./log",

	RotateUnit:   "h",
	RotateCount:  24,
	RotateSwitch: false,
}

func Init(conf LogConfig) *zap.SugaredLogger {
	logConfig.ZapLevel = getLogLevel(conf.Level)
	logConfig.Log2File = true
	logConfig.Path = env.GetLogDirPath()

	if env.IDC == env.IDC_TEST {
		// 开发环境下默认输出到文件，支持自定义是否输出到终端
		logConfig.Stdout = conf.Stdout
	} else {
		// 线上
		logConfig.Stdout = false
	}

	if _, err := os.Stat(logConfig.Path); os.IsNotExist(err) {
		err = os.MkdirAll(logConfig.Path, 0777)
		if err != nil {
			panic(fmt.Errorf("log conf err: create log dir '%s' error: %s", logConfig.Path, err))
		}
	}

	flagFile := path.Join(logConfig.Path, ".rotate")
	if conf.Rotate.Switch {
		conf.Rotate.Unit = strings.ToUpper(conf.Rotate.Unit)
		if !rotateUnitValid(conf.Rotate.Unit) {
			panic("rotate unit only support 天(D/d)、小时(H/h)、分钟M(M/m)")
		}
		logConfig.RotateSwitch = true
		logConfig.RotateUnit = conf.Rotate.Unit
		logConfig.RotateCount = conf.Rotate.Count

		// 在日志目录增加一个文件表示使用框架切割
		if fd, err := os.Create(flagFile); err != nil {
			panic(".rotate file create error: " + err.Error())
		} else {
			_, _ = fd.WriteString(time.Now().String())
			_ = fd.Close()
		}
	} else {
		if _, err := os.Stat(flagFile); !os.IsNotExist(err) {
			if err := os.Remove(flagFile); err != nil {
				panic(".rotate file remove error: " + err.Error())
			}
		}
	}

	ServerLogger = GetLogger()
	return ServerLogger
}

func rotateUnitValid(when string) bool {
	switch when {
	case "M", "H", "D":
		return true
	default:
		return false
	}
}
