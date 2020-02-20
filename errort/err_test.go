package errort

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorf(t *testing.T) {
	Init("C:/Users/zrufo/go/src/github.com/ZR233/goutils")
	errT := errors.New("test error")
	err := warpedErr(errT)
	err = fmt.Errorf("wrap error2%w", err)
	t.Log(GetTrace(err))

	if errors.Is(err, errT) {
		t.Log("ok")
	}
}
