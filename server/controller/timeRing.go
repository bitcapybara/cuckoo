package controller

import (
	"github.com/bitcapybara/cuckoo/core"
	"sync"
)

type timeRing struct {
	ringData map[int][]core.Job
	mu       sync.Mutex
}

func NewTimeRing() *timeRing {
	return &timeRing{
		ringData: make(map[int][]core.Job),
	}
}

func (t *timeRing) getAndRemove(second int) []core.Job {
	t.mu.Lock()
	defer t.mu.Unlock()
	jobs := t.ringData[second]
	delete(t.ringData, second)
	return jobs
}

func (t *timeRing) put(second int, job core.Job) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if jobs, ok := t.ringData[second]; !ok {
		t.ringData[second] = []core.Job{job}
	} else {
		jobs = append(jobs, job)
	}
}
