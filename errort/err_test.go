package errort

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorf(t *testing.T) {

	errT := errors.New("test error")
	err := warpedErr(errT)
	err = fmt.Errorf("wrap error2%w", err)
	t.Log(GetTrace(err))

	if errors.Is(err, errT) {
		t.Log("ok")
	}
}

func Test_getFirstErrorStd(t *testing.T) {
	err := errors.New("test1")
	err = fmt.Errorf("test2%w", err)
	err = NewFromError(err, 0)
	err = fmt.Errorf("test3%w", err)
	//err2 := GetStdError(err)
	println(err.Error())
	//println(err2.Error())
}
