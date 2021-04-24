package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type FixedRateScheduler struct {
}

func newScheduleFixedRate(rule core.ScheduleRule) FixedRateScheduler {
	return FixedRateScheduler{}
}

func (f FixedRateScheduler) Next(time time.Time) time.Time {
	panic("implement me")
}
