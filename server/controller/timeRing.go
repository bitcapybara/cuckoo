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
	jobIds := t.ringData[second]
	delete(t.ringData, second)
	return jobIds
}

func (t *timeRing) put(second int, job core.Job) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if ids, ok := t.ringData[second]; !ok {
		t.ringData[second] = []core.Job{job}
	} else {
		ids = append(ids, job)
	}
}

func (t *timeRing) Replace(second int, jobs []core.Job) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ringData[second] = jobs
}
