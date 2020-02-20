/*
@Time : 2019-11-11 13:43
@Author : zr
*/
package errort

import (
	"errors"
	"runtime/debug"
	"strings"
)

type ErrorWithTrace interface {
	error
	Unwrap() error
	ErrorWithTrace() string
}

type StdError struct {
	err          error
	msg          string
	msgWithTrace string
}

func (s *StdError) Error() string {
	return s.msg
}
func (s *StdError) ErrorWithTrace() string {
	return s.msgWithTrace
}
func (s *StdError) Unwrap() error {
	return s.err
}

func NewFromError(err error, warpCount int) *StdError {
	stdErr := &StdError{
		err:          err,
		msgWithTrace: fileLocation(warpCount + 1),
		msg:          err.Error(),
	}
	return stdErr
}
func NewErr(msg string, warpCount int) *StdError {
	err := errors.New(msg)
	stdErr := &StdError{
		err:          err,
		msgWithTrace: fileLocation(warpCount + 1),
		msg:          msg + "\n" + err.Error(),
	}
	return stdErr
}

func GetStdError(err error) (stdErr *StdError) {
	var errs []*StdError
	for {
		if err == nil {
			break
		}
		ok := false
		if stdErr, ok = err.(*StdError); ok {
			errs = append(errs, stdErr)
		}
		err = errors.Unwrap(err)
	}
	l := len(errs)
	if l > 0 {
		stdErr = errs[l-1]
	}
	return
}

func GetTrace(err error) (trace string) {
	if err == nil {
		return
	}
	stdErr := GetStdError(err)
	if stdErr != nil {
		trace = stdErr.ErrorWithTrace()
	}
	return
}

//生成一个错误
// @warpCount 错误出现位置忽略层级。用于被封装后，在调用堆栈中忽略封装函数名
func fileLocation(warpCount int) string {
	trace := string(debug.Stack())
	traceArr := strings.Split(trace, "\n")
	level := 5 + warpCount*2
	traceArr = traceArr[level:]

	return strings.Join(traceArr, "\n")
}
