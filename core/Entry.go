package core

import (
	"github.com/bitcapybara/cuckoo/core/schedule"
	"time"
)

type Entry struct {
	ID string
	Schedule schedule.Schedule
	Next time.Time
	Prev time.Time
	Job Job
}
