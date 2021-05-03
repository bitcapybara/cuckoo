package schedule

import (
	"errors"
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type Scheduler interface {
	Next(time.Time) time.Time
}

func GetScheduler(rule core.ScheduleRule) (scheduler Scheduler, err error) {
	switch rule.ScheduleType {
	case core.Cron:
		if rule.ParseOption == 0 {
			s, pErr := ParseStandard(rule.CronExpr)
			if pErr != nil {
				return nil, pErr
			} else {
				scheduler = s
			}
		} else {
			s, pErr := NewParser(rule.ParseOption).Parse(rule.CronExpr)
			if pErr != nil {
				return nil, pErr
			} else {
				scheduler = s
			}
		}
	case core.FixedRate:
		scheduler = newScheduleFixedRate(rule)
	default:
		return nil, errors.New("不支持的调度类型")
	}
	return
}

// 调度策略是无状态的，每次调用时生成新对象
func Schedule(rule core.ScheduleRule, prev time.Time) time.Time {
	scheduler, err := GetScheduler(rule)
	if err != nil {
		return time.Time{}
	}
	return scheduler.Next(prev)
}
