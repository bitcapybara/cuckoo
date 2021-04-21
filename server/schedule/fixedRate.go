package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type scheduleFixedRate struct {
}

func newScheduleFixedRate(rule core.ScheduleRule) scheduleFixedRate {
	return scheduleFixedRate{}
}

func (f scheduleFixedRate) Next(time time.Time) time.Time {
	panic("implement me")
}
