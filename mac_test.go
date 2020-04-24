package goutils

import (
	"testing"
)

func TestMacStringToBytes(t *testing.T) {
	mac := "03:ac:ef:13:34:1f"
	b, e := MacStringToBytes(mac, ":")
	if e != nil {
		t.Error(e)
	}

	mac2, e := MacStringFromBytes(b, ":")

	if mac != mac2 {
		t.Error("")
	}

}
