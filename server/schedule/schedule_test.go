package schedule

import (
	"github.com/bitcapybara/cuckoo/core"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	for i := 0; i < 200; i++ {
		time.Sleep(time.Millisecond * 800)
		now := time.Now()
		scheduleTime := Schedule(core.ScheduleRule{
			ScheduleType: core.Cron,
			CronExpr:     "0/5 * * * *",
			ParseOption:  core.Second | core.Minute | core.Hour | core.Dom | core.Month,
		}, now)
		println(now.Format("2006-01-02 15:04:05.000") + "------" + scheduleTime.Format("2006-01-02 15:04:05.000"))
	}
}
