/*
@Time : 2019-10-10 14:08
@Author : zr
*/
package errors

import "golang.org/x/xerrors"

//可用函数
var F = Utils{}

type Type int

var typeMessageMap_ map[Type]string

type Error struct {
	errType Type
	trace   error
}

const UnknownType Type = -1

//必须在最开始设置
func SetErrorTypeMessageMap(typeMessageMap map[Type]string) {
	typeMessageMap_ = typeMessageMap
}

func (e *Error) Error() string {
	if v, ok := typeMessageMap_[e.errType]; ok {
		return v
	} else {
		return "error not defined"
	}
}
func (e *Error) Trace() string {
	return e.trace.Error()
}
func (e *Error) Code() int {
	return int(e.errType)
}

type Utils struct {
}

func (u Utils) Wrap(err error) *Error {
	if err_, ok := err.(*Error); ok {
		err_.trace = xerrors.Errorf("%w", err_.trace)
		return err_
	} else {
		e := F.New(UnknownType)
		return e
	}
}

func (u Utils) New(errType Type) *Error {
	err := &Error{}
	err.errType = errType
	err.trace = xerrors.New(err.Error())
	return err
}

func (u Utils) Equal(err error, errType Type) bool {
	if err_, ok := err.(*Error); ok {
		return err_.errType == errType
	}
	return false
}
