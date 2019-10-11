/*
@Time : 2019-10-10 14:08
@Author : zr
*/
package errors

import (
	"runtime/debug"
)

//可用函数
var F = Utils{}

type Type int

var typeMessageMap_ map[Type]string

type Error struct {
	message string
	errType Type
	trace   string
}

const UnknownType Type = -1

//必须在最开始设置
func SetErrorTypeMessageMap(typeMessageMap map[Type]string) {
	typeMessageMap_ = typeMessageMap
}

func (e *Error) Error() string {
	if e.message != "" {
		return e.message
	}

	if v, ok := typeMessageMap_[e.errType]; ok {
		return v
	} else {
		return "error not defined"
	}
}
func (e *Error) Trace() string {
	return e.trace
}
func (e *Error) Code() int {
	return int(e.errType)
}

type Utils struct {
}

func (u Utils) Wrap(err error) *Error {
	if err_, ok := err.(*Error); ok {
		return err_
	} else {
		e := F.New(UnknownType, err.Error())
		err_.trace = string(debug.Stack())
		return e
	}
}

func (u Utils) New(errType Type, msg string) *Error {
	err := &Error{}
	err.errType = errType
	if msg == "" {
		msg = err.Error()
	} else {
		err.message = msg
	}
	err.trace = string(debug.Stack())
	//err.trace = xerrors.New(msg)
	return err
}

func (u Utils) Equal(err error, errType Type) bool {
	if err_, ok := err.(*Error); ok {
		return err_.errType == errType
	}
	return false
}
