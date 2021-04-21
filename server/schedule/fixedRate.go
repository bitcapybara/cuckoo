package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type scheduleFixedRate struct {
}

func newScheduleFixedRate() scheduleFixedRate {
	return scheduleFixedRate{}
}

func (f scheduleFixedRate) Next(rule core.ScheduleRule, time time.Time) time.Time {
	panic("implement me")
}
