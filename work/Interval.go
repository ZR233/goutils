package work

import "time"

//IntervalWork 间隔时间任务
type IntervalWork struct {
	intervalTime time.Duration
	workBase
}

func (i *IntervalWork) loop() {
	defer loopDefer(i.onError)
	err := i.loopFunc(i.ctx)
	if err != nil {
		i.onError(err)
	}
}

func (i *IntervalWork) Run() {

	for {
		i.loop()
		select {
		case <-i.ctx.Done():
			return
		case <-time.After(i.intervalTime):

		}
	}
}
