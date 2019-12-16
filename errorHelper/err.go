package errorHelper

import (
	"errors"
	"runtime/debug"
	"strings"
)

type ErrorWithTrace struct {
	error
	trace string
}

func (e *ErrorWithTrace) Error() string {
	return e.error.Error()
}

func (e *ErrorWithTrace) Unwrap() error {
	return e.error
}

func WithTrace(err error) *ErrorWithTrace {
	e := &ErrorWithTrace{}
	e.error = err

	trace := string(debug.Stack())
	traceArr := strings.Split(trace, "\n")
	level := 5

	traceArr_ := append(traceArr[:1], traceArr[level:]...)
	e.trace = strings.Join(traceArr_, "\n")

	return e
}

func Trace(err error) string {
	for err != nil {
		if e, ok := err.(*ErrorWithTrace); ok {
			return e.trace
		}
		err = errors.Unwrap(err)
	}
	return ""
}
