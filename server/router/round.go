package router

import (
	"github.com/bitcapybara/cuckoo/core"
	"sync"
)

type routeRound struct {
	lastAddrIndex map[string]int
	mu sync.Mutex
}

func newRouteRound() *routeRound {
	return &routeRound{
		lastAddrIndex: make(map[string]int),
	}
}

func (r *routeRound) route(groupName string, clients []core.NodeAddr) core.NodeAddr {
	r.mu.Lock()
	defer r.mu.Unlock()
	index := 0
	if lastIndex, ok := r.lastAddrIndex[groupName]; ok {
		index = (lastIndex + 1) % len(clients)
	}
	r.lastAddrIndex[groupName] = index
	return clients[index]
}
