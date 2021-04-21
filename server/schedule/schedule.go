package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type Scheduler interface {
	Next(core.ScheduleRule, time.Time) time.Time
}

var schedulers = map[core.ScheduleType]Scheduler{
	core.Cron: newScheduleCron(),
	core.FixedDelay: newScheduleFixedDelay(),
	core.FixedRate: newScheduleFixedRate(),
}

func GetScheduler(scheduleType core.ScheduleType) Scheduler {
	return schedulers[scheduleType]
}
