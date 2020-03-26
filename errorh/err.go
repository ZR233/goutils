package errorh

import (
	"bytes"
	"errors"
	"runtime/debug"
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

	e.trace = trace()

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

func trace() string {
	s := []byte("errorh/err.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := debug.Stack()
	start := -1
	for {
		start = bytes.Index(stack, s)
		if start < 0 {
			break
		}
		stack = stack[start+len(s):]
	}

	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return string(stack)
}
