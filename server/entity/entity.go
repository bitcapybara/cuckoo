package entity

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type JobInfo struct {
	Job    core.Job
	Enable bool
	Next   time.Time
	Prev   time.Time
}
