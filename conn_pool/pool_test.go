package conn_pool

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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
	id     int64
	closed bool
}

func (t *testCloser) Close() error {
	fmt.Printf("%d close", t.id)
	t.closed = true
	return nil
}

var id = int64(0)

func testFactory() ConnFactory {
	return func() (closer io.Closer, e error) {
		c := &testCloser{}
		c.id = atomic.AddInt64(&id, 1)
		c.closed = false
		closer = c
		return
	}
}
func testConnFailFactory() ConnFactory {
	return func() (closer io.Closer, e error) {
		c := &testCloser{}
		c.id = atomic.AddInt64(&id, 1)
		c.closed = true
		closer = c
		return
	}
}

func testErrHandler() ErrorHandler {
	return func(err error) {
		println(err)
	}
}

func testConnTestFunc() ConnTestFunc {
	return func(closer io.Closer) bool {
		return !closer.(*testCloser).closed
	}
}

func TestPool_NewPool(t *testing.T) {
	pool, _ := NewPool(testFactory(), testErrHandler(), testConnTestFunc())
	defer pool.Close()
	var (
		conn1 io.Closer
		conn2 io.Closer
		conn3 io.Closer
		err   error
	)

	conn1, err = pool.Acquire()
	if err != nil {
		t.Error(err)
	}

	if pool.numOpen != 1 || len(pool.queue) != 0 || len(pool.pool) != 1 {
		t.Error(fmt.Sprintf("获取一个连接:numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
		return
	}

	pool.Release(conn1)

	if pool.numOpen != 1 || len(pool.queue) != 1 || len(pool.pool) != 1 {
		t.Error(fmt.Sprintf("释放一个连接:numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
		return
	}

	_ = pool.Close()
	<-time.After(time.Second)
	if pool.numOpen != 0 || len(pool.queue) != 0 || len(pool.pool) != 0 {
		t.Error(fmt.Sprintf("关闭: numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
	}

	pool, _ = NewPool(testFactory(), testErrHandler(), testConnTestFunc(), OptionMaxOpen(3))
	//<-time.After(time.Millisecond * 200)
	//if pool.numOpen != 1 || len(pool.queue) != 1 || len(pool.pool) != 1 {
	//	t.Error(fmt.Sprintf("获取2个连接池: numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
	//	return
	//}

	conn1, err = pool.Acquire()
	if err != nil {
		t.Error(err)
		return
	}
	if pool.numOpen != 2 || len(pool.queue) != 1 || len(pool.pool) != 2 {
		t.Error(fmt.Sprintf("获取第1个连接 numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
		return
	}

	conn2, err = pool.Acquire()
	if err != nil {
		t.Error(err)
		return
	}

	if pool.numOpen != 2 || len(pool.queue) != 0 || len(pool.pool) != 2 {
		t.Error(fmt.Sprintf("获取第2个连接 numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
		return
	}

	pool.Release(conn1)
	conn1.(*testCloser).closed = true
	if pool.numOpen != 2 || len(pool.queue) != 1 || len(pool.pool) != 2 {
		t.Error(fmt.Sprintf("释放一个 numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
		return
	}

	pool.CloseOne(conn2)

	if pool.numOpen != 1 || len(pool.queue) != 1 || len(pool.pool) != 1 {
		t.Error(fmt.Sprintf("关闭一个 numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
		return
	}
	conn1, err = pool.Acquire()
	conn2, err = pool.Acquire()
	conn3, err = pool.Acquire()
	if err != nil {
		t.Error(err)
		return
	}

	if pool.numOpen != 3 || len(pool.queue) != 0 || len(pool.pool) != 3 {
		t.Error(fmt.Sprintf("再获取可用连接 numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
		return
	}
	if conn1.(*testCloser).closed {
		t.Error("获取已关闭连接")
	}
	if conn2.(*testCloser).closed {
		t.Error("获取已关闭连接")
	}
	if conn3.(*testCloser).closed {
		t.Error("获取已关闭连接")
	}
}

func TestPool_Pressure(t *testing.T) {
	pool, _ := NewPool(testFactory(), testErrHandler(), testConnTestFunc(), OptionMaxOpen(10))

	thread := 2000

	wg := sync.WaitGroup{}

	for i := 0; i < thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := pool.Acquire()
			if err != nil {
				t.Error(err)
			}
			defer pool.Release(conn)
			n := time.Duration(rand.Intn(10))

			<-time.After(time.Millisecond * (20 + n))

		}()
	}
	wg.Wait()
	if pool.numOpen != 10 || len(pool.queue) != 10 || len(pool.pool) != 10 {
		t.Error(fmt.Sprintf("numOpen %d queue len %d pool len %d ", pool.numOpen, len(pool.queue), len(pool.pool)))
	}

}

func TestPool_Acquire(t *testing.T) {
	pool, _ := NewPool(testConnFailFactory(), testErrHandler(), testConnTestFunc(), OptionMaxOpen(2))

	tests := []struct {
		name    string
		pool    *Pool
		want    io.Closer
		wantErr bool
	}{
		{"从非正常服务器获取连接", pool, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.pool.Acquire()
			if (err != nil) != tt.wantErr {
				t.Errorf("Acquire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Acquire() got = %v, want %v", got, tt.want)
			}
		})
	}
}
