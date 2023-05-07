package base

import (
	"fmt"
	"github.com/pkg/errors"
)

type Error struct {
	ErrNo  int
	ErrMsg string
}

type logger interface {
	Print(v ...interface{})
}

func NewError(code int, message, userMsg string) Error {
	return Error{
		ErrNo:  code,
		ErrMsg: message,
	}
}

func (err Error) Error() string {
	return err.ErrMsg
}

func (err Error) Sprintf(v ...interface{}) Error {
	err.ErrMsg = fmt.Sprintf(err.ErrMsg, v...)
	return err
}

func (err Error) Equal(e error) bool {
	switch errors.Cause(err).(type) {
	case Error:
		return err.ErrNo == errors.Cause(err).(Error).ErrNo
	default:
		return false
	}
}
