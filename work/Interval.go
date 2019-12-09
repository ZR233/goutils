package work

import "time"

//IntervalWork 间隔时间任务
type IntervalWork struct {
	intervalTime time.Duration
	workBase
}

func (i *IntervalWork) loop() {
	defer loopDefer(i.onError)
	i.workHandler()
}

func (i *IntervalWork) Run() {
	for {
		select {
		case <-i.ctx.Done():
			return
		case <-time.After(i.intervalTime):
			i.loop()
		}
	}
}
func (i *IntervalWork) Stop() {
	i.cancel()
}
