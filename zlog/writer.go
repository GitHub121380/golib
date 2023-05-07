package zlog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/GitHub121380/file-rotatelogs"
)

// 带软链形式的日志切割方式
func getRotateWriter(filename string) (wr io.Writer, err error) {
	r := newRotate(logConfig.RotateUnit, uint(logConfig.RotateCount))

	wr, err = rotatelogs.New(
		filename,
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithRotationRegRule(r.getRotateRegRule()),
		rotatelogs.WithRotationCount(r.getRotateCount()),
		rotatelogs.WithRotationTime(r.getRotationTime()),
	)

	return wr, err
}

func getFileWriter(filename string) (wr io.Writer, err error) {
	logDir := logConfig.Path
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0777)
		if err != nil {
			panic(fmt.Errorf("log conf err: create log dir '%s' error: %s", logDir, err))
		}
	}

	if fi, err := os.Lstat(filename); err == nil {
		if mode := fi.Mode(); mode&os.ModeSymlink != 0 {
			_ = os.Remove(filename)
		}
	}

	wr, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	return wr, err
}

type rotate struct {
	unit         string
	count        uint
	regRule      string
	timeInterval time.Duration
}

func newRotate(unit string, count uint) *rotate {
	if unit == "" {
		unit = "H"
	}

	if count <= 0 {
		count = 24
	}

	// 底层rotate是保留的日志文件的总数，与上层业务理解不一致
	// 上层count表示备份文件的个数
	count += 1

	var reg string
	var rotateTime time.Duration
	switch strings.ToUpper(unit) {
	case "D": // day
		reg = "%Y%m%d"
		rotateTime = 24 * time.Hour
	case "M": // minute
		reg = "%Y%m%d%H%M"
		rotateTime = 1 * time.Minute
	case "H": // hour, by default
		fallthrough
	default:
		reg = "%Y%m%d%H"
		rotateTime = 1 * time.Hour
	}

	return &rotate{
		unit:         unit,
		count:        count,
		regRule:      reg,
		timeInterval: rotateTime,
	}
}

func (r *rotate) getRotateUnit() string {
	return r.unit
}

func (r *rotate) getRotateCount() uint {
	return r.count
}

func (r *rotate) getRotationTime() time.Duration {
	return r.timeInterval
}

func (r *rotate) getRotateRegRule() string {
	return r.regRule
}
