package goutils

import "testing"

func TestIDGen(t *testing.T) {
	sid := uint32(0)
	for i := 0; i < 500; i++ {
		id := IDGen(0, &sid)
		println(id)
	}
}
