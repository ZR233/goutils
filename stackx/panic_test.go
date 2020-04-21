package stackx

import (
	"testing"
)

func hello(str string) {
	panic("test p")
}

func TestRecover(t *testing.T) {
	defer Recover()

	hello("param 1")

	println("awd")

}
