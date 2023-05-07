package zlog

import (
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 业务日志输出logger
var (
	ServerLogger *zap.SugaredLogger
	ModuleLogger *zap.Logger
	AccessLogger *zap.Logger
)

const (
	LogNameAccess = "access"
	LogNameModule = "module"
	LogNameServer = "server"
)

const (
	txtLogNormal    = "normal"
	txtLogWarnFatal = "warnfatal"
	txtLogStdout    = "stdout"
)

// NewLogger 新建Logger，每一次新建会同时创建x.log与x.log.wf (access.log 不会生成wf)
func newLogger(name string) *zap.Logger {

	var infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl <= zapcore.InfoLevel
	})

	var errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl >= zapcore.WarnLevel
	})

	var stdLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl >= zapcore.DebugLevel
	})

	var zapCore []zapcore.Core
	if logConfig.Log2File {
		c := zapcore.NewCore(
			getEncoder(),
			zapcore.AddSync(getLogWriter(name, txtLogNormal)),
			infoLevel)
		zapCore = append(zapCore, c)

		if name != LogNameAccess {
			c := zapcore.NewCore(
				getEncoder(),
				zapcore.AddSync(getLogWriter(name, txtLogWarnFatal)),
				errorLevel)
			zapCore = append(zapCore, c)
		}
	}

	if logConfig.Stdout {
		c := zapcore.NewCore(
			getEncoder(),
			zapcore.AddSync(getLogWriter(name, txtLogStdout)),
			stdLevel)
		zapCore = append(zapCore, c)
	}

	// core
	core := zapcore.NewTee(zapCore...)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	// 由于之前没有DPanic，同化DPanic和Panic
	development := zap.Development()

	// 设置初始化字段
	filed := zap.Fields()

	// 构造日志
	logger := zap.New(core, filed, caller, development)

	return logger
}

func getLogLevel(lv string) (level zapcore.Level) {
	str := strings.ToUpper(lv)
	switch str {
	case "DEBUG":
		level = zap.DebugLevel
	case "INFO":
		level = zap.InfoLevel
	case "WARN":
		level = zap.WarnLevel
	case "ERROR":
		level = zap.ErrorLevel
	case "FATAL":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return level
}

func getEncoder() zapcore.Encoder {
	// 公用编码器
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	encoderCfg := zapcore.EncoderConfig{
		LevelKey:      "level",
		TimeKey:       "time",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//LineEnding:    "tp=&tc=xxx.logger\n",
		EncodeCaller: zapcore.FullCallerEncoder, // 全路径编码器
		//EncodeName:     zapcore.FullNameEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewJSONEncoder(encoderCfg)
}

func getLogWriter(name, logType string) (wr io.Writer) {
	// stdOut
	if logType == txtLogStdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}

	// 写日志到文件filename中
	filename := genFilename(name, logType)
	return NewTimeFileLogWriter(filename)
}

func genFilename(appName, logType string) string {
	var tailFixed string
	switch logType {
	case txtLogNormal:
		tailFixed = ".log"
	case txtLogWarnFatal:
		tailFixed = ".log.wf"
	default:
		tailFixed = ".log"
	}

	return appName + tailFixed
}

func CloseLogger() {
	if ServerLogger != nil {
		_ = ServerLogger.Sync()
	}

	if ModuleLogger != nil {
		_ = ModuleLogger.Sync()
	}

	if AccessLogger != nil {
		_ = AccessLogger.Sync()
	}
}

// 避免用户改动过大，以下为封装的之前的Entry打印field的方法
type Fields map[string]interface{}
type entry struct {
	s *zap.SugaredLogger
}

func NewEntry(s *zap.SugaredLogger) *entry {
	x := s.Desugar().WithOptions(zap.AddCallerSkip(+1)).Sugar()
	return &entry{s: x}
}

// 注意这种使用方式固定头的顺序会变
func (e entry) WithFields(f Fields) *zap.SugaredLogger {
	var fields []interface{}
	for k, v := range f {
		fields = append(fields, k, v)
	}

	return e.s.With(fields...)
}
