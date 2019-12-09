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
}

type OnError func(err error)

type workBase struct {
	workHandler func()
	onError     OnError
	ctx         context.Context
	cancel      context.CancelFunc
	stopped     chan bool
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

func initWorkBase(workBase *workBase, workHandler func(), onError OnError) {
	workBase.stopped = make(chan bool)
	workBase.onError = onError
	workBase.ctx, workBase.cancel = context.WithCancel(context.Background())
	workBase.workHandler = workHandler
}

func RegisterNewWork(intervalTime time.Duration, workHandler func(), onError OnError) *IntervalWork {
	work := &IntervalWork{
		intervalTime: intervalTime,
	}
	initWorkBase(&work.workBase, workHandler, onError)

	go work.Run()
	return work
}

func RegisterScheduleDayWork(schedule time.Time, workHandler func(), onError OnError) *ScheduleDayWork {

	work := &ScheduleDayWork{
		ScheduleTime: schedule,
	}
	initWorkBase(&work.workBase, workHandler, onError)
	go work.Run()
	return work
}
