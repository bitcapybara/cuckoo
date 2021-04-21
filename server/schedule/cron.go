package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type scheduleCron struct {
}

func newScheduleCron(rule core.ScheduleRule) scheduleCron {
	return scheduleCron{}
}

func (s scheduleCron) Next(time time.Time) time.Time {
	panic("implement me")
}
