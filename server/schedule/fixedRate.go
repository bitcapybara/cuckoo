package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type FixedRateScheduler struct {
	Init     time.Duration
	Duration time.Duration
}

func newScheduleFixedRate(rule core.ScheduleRule) FixedRateScheduler {
	init := rule.Initial
	duration := rule.Duration
	if duration < time.Second {
		duration = time.Second
	}
	if init < time.Second {
		init = time.Second
	}
	return FixedRateScheduler{
		Init:     init - time.Duration(duration.Nanoseconds())%time.Second,
		Duration: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

func (schedule FixedRateScheduler) Next(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().Add(schedule.Init - time.Duration(t.Nanosecond())*time.Nanosecond)
	}
	return t.Add(schedule.Duration - time.Duration(t.Nanosecond())*time.Nanosecond)
}
