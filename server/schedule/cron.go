package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type scheduleCron struct {
}

func newScheduleCron() scheduleCron {
	return scheduleCron{}
}

func (s scheduleCron) Next(rule core.ScheduleRule, time time.Time) time.Time {
	panic("implement me")
}
