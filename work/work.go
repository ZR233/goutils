/*
@Time : 2019-12-09 09:20
@Author : zr
*/
package work

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Work interface {
	Run()
	Stop()
	Join()
	getWorkBase() *workBase
}

type OnError func(err error)
type LoopFunc func(ctx context.Context) (err error)

type workBase struct {
	loopFunc LoopFunc
	onError  OnError
	ctx      context.Context
	cancel   context.CancelFunc
	stopped  chan bool
}

func loopDefer(onError OnError) {
	if p := recover(); p != nil {
		if err, ok := p.(error); ok {
			onError(err)
		} else {
			onError(errors.New(fmt.Sprint("work error:\n", p)))
		}
	}
}

func initWork(work Work, loopFunc LoopFunc, onError OnError) {
	workBase := work.getWorkBase()
	workBase.stopped = make(chan bool)
	workBase.onError = onError
	workBase.ctx, workBase.cancel = context.WithCancel(context.Background())
	workBase.loopFunc = loopFunc

	go func() {
		defer func() {
			workBase.stopped <- true
		}()
		work.Run()
	}()
}

func RegisterNewWork(intervalTime time.Duration, loopFunc LoopFunc, onError OnError) *IntervalWork {
	work := &IntervalWork{
		intervalTime: intervalTime,
	}
	initWork(work, loopFunc, onError)

	return work
}

func RegisterScheduleDayWork(schedule time.Time, loopFunc LoopFunc, onError OnError) *ScheduleDayWork {

	work := &ScheduleDayWork{
		ScheduleTime: schedule,
	}
	initWork(work, loopFunc, onError)

	return work
}

func (w *workBase) getWorkBase() *workBase {
	return w
}

func (w *workBase) Stop() {
	w.cancel()
}
func (w *workBase) Join() {
	<-w.stopped
}
