/*
@Time : 2019-11-11 13:47
@Author : zr
*/
package goutils

import (
	"github.com/ZR233/goutils/errors2"
	"strings"
)

import "testing"

var (
	testError1 = errors2.U.NewErrorType(1, "测试error1 debug", "测试error1")
	testError2 = errors2.U.NewErrorType(2, "测试error2 debug", "测试error2")
)

func throughError() *errors2.StdError {
	return errors2.U.NewStdErrorWarpInFunc(testError1, "throughError", 1)
}
func FuncWithError() *errors2.StdError {
	return throughError()
}

func TestNewError(t *testing.T) {
	type args struct {
		error *errors2.StdError
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		//{"", args{throughError()}, ""},
		{"", args{FuncWithError()}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trace := tt.args.error.Trace()

			arr := strings.Split(trace, "\n")

			if arr[1] != "github.com/ZR233/goutils.FuncWithError(0x0)" {
				t.Error("不相等")
			}
		})
	}
}
