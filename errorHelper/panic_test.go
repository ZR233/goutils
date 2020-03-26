package errorHelper

import (
	"testing"
)

func TestPanicTrace(t *testing.T) {
	defer println(PanicTrace(4))

	panic("123")

}
