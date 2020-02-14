/*
@Time : 2019-11-11 13:43
@Author : zr
*/
package erroru

import (
	"fmt"
	"runtime/debug"
	"strings"
)

var srcPrefix string

func Init(srcRootPath string) (err error) {
	srcPrefix = srcRootPath
	return
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

func Warp(err error) *StdError {
	return AddInfo("", err)
}

func AddInfo(msg string, err error) *StdError {
	stdErr := &StdError{
		err:          err,
		msgWithTrace: fmt.Sprintln(fileLocation(2)+"\t|"+msg) + StdErrorFrom(err).ErrorWithTrace(),
		msg:          msg,
	}
	return stdErr
}
func StdErrorFrom(err error) (stdErr *StdError) {
	ok := false
	if stdErr, ok = err.(*StdError); !ok {
		stdErr = &StdError{
			err:          err,
			msg:          err.Error(),
			msgWithTrace: err.Error(),
		}
	}
	return
}

//生成一个错误
// @warpCount 错误出现位置忽略层级。用于被封装后，在调用堆栈中忽略封装函数名
func fileLocation(warpCount int) string {
	trace := string(debug.Stack())
	traceArr := strings.Split(trace, "\n")
	level := 4 + warpCount*2
	l := strings.TrimLeft(traceArr[level], "\t")
	l = strings.Replace(l, srcPrefix, "", 1)
	return l
}
