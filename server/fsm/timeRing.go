package fsm

import (
	"github.com/bitcapybara/cuckoo/core"
	"sync"
)

type TimeRing struct {
	ringData map[int][]core.JobId
	mu       sync.Mutex
}

func NewTimeRing() *TimeRing {
	return &TimeRing{
		ringData: make(map[int][]core.JobId),
	}
}

func (t *TimeRing) Get(second int) []core.JobId {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.ringData[second]
}

func (t *TimeRing) GetAndRemove(second int) []core.JobId {
	t.mu.Lock()
	defer t.mu.Unlock()
	jobIds := t.ringData[second]
	delete(t.ringData, second)
	return jobIds
}

func (t *TimeRing) Put(second int, id core.JobId) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if ids, ok := t.ringData[second]; !ok {
		t.ringData[second] = []core.JobId{id}
	} else {
		ids = append(ids, id)
	}
}

func (t *TimeRing) Replace(second int, jobIds []core.JobId) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ringData[second] = jobIds
}
