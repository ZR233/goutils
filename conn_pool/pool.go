package conn_pool

import (
	"context"
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

// 创建连接的方法
type ConnFactory func() (io.Closer, error)

// 错误处理
type ErrorHandler func(err error)

// 判断连接是否可用
type ConnTestFunc func(closer io.Closer) bool

type Pool struct {
	sync.Mutex
	pool       map[io.Closer]*connStatus
	queue      chan io.Closer
	numOpen    int // 当前池中资源数
	ctx        context.Context
	cancel     context.CancelFunc
	config     *config
	createConn chan bool
	stop       bool
}

type connStatus struct {
	createTime time.Time
	using      bool
}

type config struct {
	maxOpen             int           // 池中最大资源数
	minOpen             int           // 池中最少资源数
	factory             ConnFactory   // 创建连接的方法
	errorHandler        ErrorHandler  // 错误处理
	connMaxAliveTime    time.Duration // 连接最大存活时间
	getConnWaitDeadline time.Duration // 获取连接最大等待时间
	connTestFunc        ConnTestFunc  // 判断连接是否可用
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

func NewPool(factory ConnFactory, errorHandler ErrorHandler, connTestFunc ConnTestFunc, options ...Option) (*Pool, error) {
	cfg := &config{
		maxOpen:             1,
		minOpen:             1,
		connMaxAliveTime:    time.Hour,
		getConnWaitDeadline: time.Second * 5,
		factory:             factory,
		errorHandler:        errorHandler,
		connTestFunc:        connTestFunc,
	}

	for _, v := range options {
		v.set(cfg)
	}

	if cfg.maxOpen <= 0 || cfg.minOpen > cfg.maxOpen {
		return nil, ErrInvalidConfig
	}
	p := &Pool{
		config:     cfg,
		pool:       map[io.Closer]*connStatus{},
		queue:      make(chan io.Closer, cfg.maxOpen),
		createConn: make(chan bool, cfg.maxOpen),
		stop:       false,
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())
	go func() {
		<-p.ctx.Done()
		p.release()
	}()

	go p.createThread()

	return p, nil
}
func (p *Pool) connOK(closer io.Closer) (b bool) {
	var err error
	defer func() {
		if pa := recover(); pa != nil {
			err = fmt.Errorf("conn test panic:\n%s", pa)
		}

		if err != nil {
			b = false
			p.handleError(err)
		}
	}()

	return p.config.connTestFunc(closer)
}
func (p *Pool) handleError(err error) {
	defer func() {
		recover()
	}()
	p.config.errorHandler(err)
}

func (p *Pool) Acquire() (io.Closer, error) {
	if p.stop {
		return nil, ErrPoolClosed
	}
	select {
	case <-p.ctx.Done():
		return nil, ErrPoolClosed
	case closer := <-p.queue:
		if p.connOK(closer) {
			return closer, nil
		} else {
			p.CloseOne(closer)
		}
	default:
	}

	p.create()
	wait := time.After(p.config.getConnWaitDeadline)

	for {
		if p.stop {
			return nil, ErrPoolClosed
		}
		select {
		case <-p.ctx.Done():
			return nil, ErrPoolClosed
		case closer := <-p.queue:
			if p.connOK(closer) {
				return closer, nil
			}
			p.CloseOne(closer)
			p.create()
			<-time.After(time.Millisecond * 20)
		case <-wait:
			return nil, ErrGetConnTimeout
		}
	}
}

func (p *Pool) createLazy() {
	for {
		if p.stop {
			return
		}
		finish := false
		select {
		case <-p.ctx.Done():
			return
		default:
			finish = p.createOneLazyFinish()
		}
		if finish {
			break
		}
	}
}

func (p *Pool) createOne() {
	// 新建连接
	closer, err := p.config.factory()
	if err != nil {
		go p.config.errorHandler(err)
		select {
		case <-p.ctx.Done():
			return
		case <-time.After(time.Millisecond * 200):
			return
		}
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

func (p *Pool) createOneLazyFinish() bool {
	p.Lock()
	defer p.Unlock()

	if p.numOpen >= p.config.maxOpen {
		return true
	}

	if p.numOpen >= p.config.minOpen && len(p.queue) > 0 {
		return true
	}

	p.createOne()
	return false
}

func (p *Pool) createThread() {
	for {
		if p.stop {
			return
		}
		select {
		case <-p.ctx.Done():
			return
		case <-p.createConn:
			p.createLazy()
			break
		case <-time.After(time.Millisecond * 500):
			p.createLazy()
			break
		}
	}
}

func (p *Pool) create() {
	defer func() {
		recover()
	}()

	select {
	case p.createConn <- true:
	default:
		return
	}
}

func (p *Pool) CloseOne(closer io.Closer) {
	p.Lock()
	defer p.Unlock()
	if p.stop {
		return
	}

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
		if p.pool != nil {
			if status, ok := p.pool[closer]; ok {
				if time.Now().Sub(status.createTime) >= p.config.connMaxAliveTime {
					p.closeWithNoLock(closer)
					r = true
				}
			}
		}
	}
	return
}

// 释放单个资源到连接池
func (p *Pool) Release(closer io.Closer) {
	p.Lock()
	defer p.Unlock()
	if p.stop {
		return
	}

	if p.connExpiredWithNoLock(closer) {
		return
	}

	if _, ok := p.pool[closer]; !ok {
		return
	}

	select {
	case <-p.ctx.Done():
		return
	case p.queue <- closer:
	default:
		return
	}
}

func (p *Pool) release() {
	p.Lock()
	defer p.Unlock()
	p.stop = true
	close(p.queue)
	p.queue = nil
	for closer := range p.pool {
		p.closeWithNoLock(closer)
	}
}

// 关闭连接池，释放所有资源
func (p *Pool) Close() error {
	p.cancel()
	return nil
}
