package main

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type jobInfo struct {
	job    core.Job
	enable bool
	next   time.Time
	prev   time.Time
}
