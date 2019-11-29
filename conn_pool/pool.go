package conn_pool

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	ErrInvalidConfig  = errors.New("invalid pool config")
	ErrPoolClosed     = errors.New("pool closed")
	ErrGetConnTimeout = errors.New("get conn timeout")
	ErrConnFactory    = errors.New("ConnFactory error")
)

type ConnFactory func() (io.Closer, error)
type ErrorHandler func(err error)

type Pool struct {
	sync.Mutex
	pool    map[io.Closer]*connStatus
	queue   chan io.Closer
	numOpen int  // 当前池中资源数
	closed  bool // 池是否已关闭
	config  *config
}

type connStatus struct {
	createTime time.Time
}

type config struct {
	maxOpen             int           // 池中最大资源数
	minOpen             int           // 池中最少资源数
	factory             ConnFactory   // 创建连接的方法
	errorHandler        ErrorHandler  // 错误处理
	connMaxAliveTime    time.Duration // 连接最大存活时间
	getConnWaitDeadline time.Duration // 获取连接最大等待时间
}

type Option interface {
	set(config *config)
}

//池中最大资源数
type OptionMaxOpen int

func (o OptionMaxOpen) set(config *config) {
	config.maxOpen = int(o)
}

//池中最少资源数
type OptionMinOpen int

func (o OptionMinOpen) set(config *config) {
	config.minOpen = int(o)
}

//连接最大存活时间
type OptionConnMaxAliveTime time.Duration

func (o OptionConnMaxAliveTime) set(config *config) {
	config.connMaxAliveTime = time.Duration(o)
}

//获取连接最大等待时间
type OptionGetConnWaitDeadline time.Duration

func (o OptionGetConnWaitDeadline) set(config *config) {
	config.getConnWaitDeadline = time.Duration(o)
}

func NewPool(factory ConnFactory, errorHandler ErrorHandler, options ...Option) (*Pool, error) {
	cfg := &config{
		maxOpen:             1,
		minOpen:             1,
		connMaxAliveTime:    time.Hour,
		getConnWaitDeadline: time.Minute,
		factory:             factory,
		errorHandler:        errorHandler,
	}

	for _, v := range options {
		v.set(cfg)
	}

	if cfg.maxOpen <= 0 || cfg.minOpen > cfg.maxOpen {
		return nil, ErrInvalidConfig
	}
	p := &Pool{
		config: cfg,
		queue:  make(chan io.Closer, cfg.maxOpen),
		pool:   map[io.Closer]*connStatus{},
	}

	for i := 0; i < p.config.minOpen; i++ {
		p.create()
	}

	return p, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}

	select {
	case closer := <-p.queue:
		return closer, nil
	case <-time.After(p.config.getConnWaitDeadline):
		return nil, ErrGetConnTimeout
	default:
		go p.create()
		select {
		case closer := <-p.queue:
			return closer, nil
		case <-time.After(p.config.getConnWaitDeadline):
			return nil, ErrGetConnTimeout
		}
	}
}

func (p *Pool) create() {
	p.Lock()
	defer p.Unlock()
	if p.closed {
		return
	}

	if len(p.queue) > 0 {

	}

	if p.numOpen >= p.config.maxOpen {
		return
	}

	// 新建连接
	closer, err := p.config.factory()
	if err != nil {
		go p.config.errorHandler(err)
		return
	}
	if closer == nil {
		panic(fmt.Errorf("%w:\nfunc success but conn is nil", ErrConnFactory))
	}

	p.numOpen++

	p.pool[closer] = &connStatus{
		createTime: time.Now(),
	}

	p.queue <- closer
}

func (p *Pool) Close(closer io.Closer) {
	p.Lock()
	defer p.Unlock()
	p.closeWithNoLock(closer)
}
func (p *Pool) closeWithNoLock(closer io.Closer) {

	if closer != nil {
		_ = closer.Close()
		delete(p.pool, closer)
	}
	p.numOpen = len(p.pool)
}

func (p *Pool) connExpiredWithNoLock(closer io.Closer) (r bool) {

	if closer != nil {
		status := p.pool[closer]
		if time.Now().Sub(status.createTime) >= p.config.connMaxAliveTime {
			p.closeWithNoLock(closer)
			r = true
		}
	}
	return
}

// 释放单个资源到连接池
func (p *Pool) Release(closer io.Closer, err error) {
	p.Lock()
	defer p.Unlock()

	if err != nil {
		p.closeWithNoLock(closer)
		return
	}

	if p.connExpiredWithNoLock(closer) {
		return
	}

	if _, ok := p.pool[closer]; !ok {
		return
	}

	select {
	case p.queue <- closer:
	default:
		return
	}
}

// 关闭连接池，释放所有资源
func (p *Pool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	defer p.Unlock()
	close(p.queue)

	for closer := range p.pool {
		p.closeWithNoLock(closer)
	}
	p.closed = true
	return nil
}