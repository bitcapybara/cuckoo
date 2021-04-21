package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type scheduleFixedDelay struct {
}

func newScheduleFixedDelay() scheduleFixedDelay {
	return scheduleFixedDelay{}
}

func (f scheduleFixedDelay) Next(rule core.ScheduleRule, time time.Time) time.Time {
	panic("implement me")
}
