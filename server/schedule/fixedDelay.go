package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type scheduleFixedDelay struct {
}

func newScheduleFixedDelay(rule core.ScheduleRule) scheduleFixedDelay {
	return scheduleFixedDelay{}
}

func (f scheduleFixedDelay) Next(time time.Time) time.Time {
	panic("implement me")
}
