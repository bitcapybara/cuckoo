package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

// 代码借用自 https://github.com/robfig/cron

// FixedDelayScheduler represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// It does not support jobs more frequent than once a second.
type FixedDelayScheduler struct {
	Init     time.Duration
	Duration time.Duration
}

// Every returns a crontab Schedule that activates once every duration.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func newScheduleFixedDelay(rule core.ScheduleRule) FixedDelayScheduler {
	init := rule.Initial
	duration := rule.Duration
	if duration < time.Second {
		duration = time.Second
	}
	if init < time.Second {
		init = time.Second
	}
	return FixedDelayScheduler{
		Init:     init - time.Duration(duration.Nanoseconds())%time.Second,
		Duration: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (schedule FixedDelayScheduler) Next(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().Add(schedule.Init - time.Duration(t.Nanosecond())*time.Nanosecond)
	}
	return t.Add(schedule.Duration - time.Duration(t.Nanosecond())*time.Nanosecond)
}
