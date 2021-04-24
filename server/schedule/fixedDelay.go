package schedule

import (
	"time"
)

// 代码借用自 https://github.com/robfig/cron

// ScheduleFixedDelay represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// It does not support jobs more frequent than once a second.
type ScheduleFixedDelay struct {
	Init  time.Duration
	Delay time.Duration
}

// Every returns a crontab Schedule that activates once every duration.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func newScheduleFixedDelay(init time.Duration, duration time.Duration) ScheduleFixedDelay {
	if duration < time.Second {
		duration = time.Second
	}
	if init < time.Second {
		init = time.Second
	}
	return ScheduleFixedDelay{
		Init: init - time.Duration(duration.Nanoseconds())%time.Second,
		Delay: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (schedule ScheduleFixedDelay) Next(t time.Time) time.Time {
	return t.Add(schedule.Delay - time.Duration(t.Nanosecond())*time.Nanosecond)
}
