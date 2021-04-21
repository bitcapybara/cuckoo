package core

import (
	"github.com/bitcapybara/cuckoo/core/schedule"
	"time"
)

type Job struct {
	ID       string
	Comment  string
	Path     string
	Schedule schedule.Schedule
	Enable   bool
	Next     time.Time
	Prev     time.Time
}
