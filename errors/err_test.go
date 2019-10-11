/*
@Time : 2019-10-11 17:16
@Author : zr
*/
package errors

import (
	"testing"
)

func TestUtils_New(t *testing.T) {
	err := F.New(UnknownType, "test")
	err = F.Wrap(err)
	t.Log(err.trace)

	t.Log(err)
}
