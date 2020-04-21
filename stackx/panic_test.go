package stackx

import (
	"testing"
)

func hello(str string) {
	panic("test p")
}

func TestRecover(t *testing.T) {
	defer func() {
		str := Stack(0)
		t.Log(str)
	}()

	hello("param 1")

	println("awd")

}
