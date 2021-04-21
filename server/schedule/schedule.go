package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type Scheduler interface {
	Next(time.Time) time.Time
}

func GetScheduler(rule core.ScheduleRule) (scheduler Scheduler) {
	switch rule.ScheduleType {
	case core.Cron:
		scheduler = newScheduleCron(rule)
	case core.FixedDelay:
		scheduler = newScheduleFixedDelay(rule)
	case core.FixedRate:
		scheduler = newScheduleFixedRate(rule)
	}
	return
}

// 调度策略是无状态的，每次调用时生成新对象
func Schedule(rule core.ScheduleRule, prev time.Time) time.Time {
	return GetScheduler(rule).Next(prev)
}
