/*
@Time : 2019-11-11 13:47
@Author : zr
*/
package goutils

import (
	"github.com/ZR233/goutils/errors2"
)

import "testing"

func TestNewError(t *testing.T) {

	errorType := errors2.U.NewErrorType(0, "", "")

	type args struct {
		errorType *errors2.ErrorType
		msg       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{errorType: errorType, msg: "测试"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := errors2.U.NewStdError(tt.args.errorType, tt.args.msg)
			t.Error(err)

		})
	}

}
