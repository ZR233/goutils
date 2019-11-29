package conn_pool

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPool_Release(t *testing.T) {

	socket, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
	buf := []byte{}
	n, err := socket.Write(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("发送%d", n)
}

type testCloser struct {
	id int64
}

func (t testCloser) Close() error {
	fmt.Printf("%d close", t.id)
	return nil
}

var id = int64(0)

func testFactory() ConnFactory {
	return func() (closer io.Closer, e error) {
		c := &testCloser{}
		c.id = atomic.AddInt64(&id, 1)
		closer = c
		return
	}
}

func testErrHandler() ErrorHandler {
	return func(err error) {
		println(err)
	}
}

func TestNewPool(t *testing.T) {
	pool, _ := NewPool(testFactory(), testErrHandler())

	conn, err := pool.Acquire()
	if err != nil {
		t.Error(err)
	}

	if pool.numOpen != 1 || len(pool.queue) != 0 || len(pool.pool) != 1 {
		t.Error("")
	}

	pool.Release(conn, nil)

	if pool.numOpen != 1 || len(pool.queue) != 1 || len(pool.pool) != 1 {
		t.Error("")
	}
	pool.Shutdown()

	pool, _ = NewPool(testFactory(), testErrHandler(), OptionMaxOpen(2))

	if pool.numOpen != 1 || len(pool.queue) != 1 || len(pool.pool) != 1 {
		t.Error("")
	}

	conn1, err := pool.Acquire()
	if err != nil {
		t.Error(err)
	}
	if pool.numOpen != 1 || len(pool.queue) != 0 || len(pool.pool) != 1 {
		t.Error("")
	}

	conn2, err := pool.Acquire()
	if err != nil {
		t.Error(err)
	}

	if pool.numOpen != 2 || len(pool.queue) != 0 || len(pool.pool) != 2 {
		t.Error("")
	}

	pool.Release(conn1, errors.New("test"))

	if pool.numOpen != 1 || len(pool.queue) != 0 || len(pool.pool) != 1 {
		t.Error("")
	}

	pool.Close(conn2)

	if pool.numOpen != 0 || len(pool.queue) != 0 || len(pool.pool) != 0 {
		t.Error("")
	}

	conn1, err = pool.Acquire()
	if err != nil {
		t.Error(err)
	}

	if pool.numOpen != 1 || len(pool.queue) != 0 || len(pool.pool) != 1 {
		t.Error("")
	}

}

func TestPool_connExpired(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		pool    map[io.Closer]*connStatus
		queue   chan io.Closer
		numOpen int
		closed  bool
		config  *config
	}
	type args struct {
		closer io.Closer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantR  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pool{
				Mutex:   tt.fields.Mutex,
				pool:    tt.fields.pool,
				queue:   tt.fields.queue,
				numOpen: tt.fields.numOpen,
				closed:  tt.fields.closed,
				config:  tt.fields.config,
			}
			if gotR := p.connExpiredWithNoLock(tt.args.closer); gotR != tt.wantR {
				t.Errorf("connExpiredWithNoLock() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
