package work

import "time"

//ScheduleWork 每日固定时间任务
type ScheduleDayWork struct {
	workBase
	ScheduleTime time.Time
}

func (s *ScheduleDayWork) loop() {
	defer loopDefer(s.onError)
	err := s.loopFunc(s.ctx)
	if err != nil {
		s.onError(err)
	}
}

func (s *ScheduleDayWork) Run() {

	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)

		next = time.Date(next.Year(), next.Month(), next.Day(), s.ScheduleTime.Hour(), s.ScheduleTime.Minute(), s.ScheduleTime.Second(), s.ScheduleTime.Nanosecond(), next.Location())

		select {
		case <-s.ctx.Done():
			return
		case <-time.After(next.Sub(now)):
			s.loop()
		}
	}
}
