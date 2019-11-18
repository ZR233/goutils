/*
@Time : 2019-11-11 13:43
@Author : zr
*/
package errors2

import (
	"runtime/debug"
	"strings"
)

var (
	//功能集
	U = Utils{}
)

type Utils struct {
}

type ErrorType struct {
	code     int
	ShowMsg  string
	DebugMsg string
}

func (Utils) NewErrorType(code int, showMsg, debugMsg string) *ErrorType {
	return &ErrorType{
		code,
		showMsg,
		debugMsg,
	}
}

type StdError struct {
	*ErrorType
	trace string
}

//将error转为*StdError, 若无法转换，则通过error生成一个*StdError，类型为errorType
func (Utils) FromError(err error, errorType *ErrorType) (stdErr *StdError) {
	if err != nil {
		ok := false
		if stdErr, ok = err.(*StdError); ok {
			return
		} else {
			if errorType == nil {
				errorType = U.NewErrorType(-1, "not define", "not define")
			}
			stdErr = U.NewStdError(errorType, err.Error())
			return
		}
	}
	return
}

//生成一个错误
func (Utils) NewStdError(errorType *ErrorType, debugMsg string) *StdError {
	return U.NewStdErrorWarpInFunc(errorType, debugMsg, 0)
}

//生成一个错误
// @warpCount 错误出现位置忽略层级。用于被封装后，在调用堆栈中忽略封装函数名
func (Utils) NewStdErrorWarpInFunc(errorType *ErrorType, debugMsg string, warpCount int) *StdError {
	trace := string(debug.Stack())
	traceArr := strings.Split(trace, "\n")
	level := 5 + warpCount*2

	traceArr_ := append(traceArr[:1], traceArr[level:]...)
	trace = strings.Join(traceArr_, "\n")

	e := &StdError{
		errorType,
		trace,
	}

	if debugMsg != "" {
		e.DebugMsg = debugMsg
	}
	return e
}

func (s *StdError) Code() int {
	return s.code
}
func (s *StdError) Error() string {
	return s.ShowMsg
}
func (s *StdError) Trace() string {
	return s.trace
}

func (s *StdError) DebugMessage() string {
	return s.DebugMsg
}
func (Utils) Equal(err error, errorType *ErrorType) bool {
	if err != nil {
		if err_, ok := err.(*StdError); ok {
			return err_.code == errorType.code
		}
	}

	return false
}
func (e *ErrorType) Equal(err error) bool {
	return U.Equal(err, e)
}
