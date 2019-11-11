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
	showMsg  string
	debugMsg string
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
	msg   string
	trace string
}

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

func (Utils) NewStdError(errorType *ErrorType, msg string) *StdError {
	trace := string(debug.Stack())
	traceArr := strings.Split(trace, "\n")
	traceArr_ := append(traceArr[:1], traceArr[5:]...)

	trace = strings.Join(traceArr_, "\n")
	return &StdError{
		errorType,
		msg,
		trace,
	}
}

func (s *StdError) Code() int {
	return s.code
}
func (s *StdError) Error() string {
	return s.showMsg
}
func (s *StdError) Trace() string {
	return s.trace
}

func (s *StdError) DebugMessage() string {
	if s.msg != "" {
		return s.msg
	} else {
		return s.debugMsg
	}
}
