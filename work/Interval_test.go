package work

import (
	"context"
	"log"
	"testing"
	"time"
)

func testOnError() OnError {
	return func(err error) {
		log.Print(err)
	}
}

type testWorkHandler struct {
	iter int
}

func (t *testWorkHandler) getFunc() LoopFunc {
	return func(ctx context.Context) (err error) {
		t.iter++
		return
	}
}

func TestIntervalWork_Run(t *testing.T) {

	tests := []struct {
		name       string
		testStruct testWorkHandler
	}{
		{"", testWorkHandler{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			work := RegisterNewWork(time.Second*1, tt.testStruct.getFunc(), testOnError())
			go work.Run()
			<-time.After(time.Second * 2)
			work.Stop()
			work.Join()

			<-time.After(time.Second * 5)
			if tt.testStruct.iter != 2 {
				t.Errorf("循环次数错误want（%d）real(%d)", 2, tt.testStruct.iter)
			}
		})
	}
}
