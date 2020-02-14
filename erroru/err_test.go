package erroru

import (
	"errors"
	"testing"
)

func TestErrorf(t *testing.T) {
	Init("C:/Users/zrufo/go/src/github.com/ZR233/goutils")
	errT := errors.New("test error")
	err := warpedErr(errT)
	err = AddInfo("wrap error2", err)
	t.Log(StdErrorFrom(err).ErrorWithTrace())

	if errors.Is(err, errT) {
		t.Log("ok")
	}
}
