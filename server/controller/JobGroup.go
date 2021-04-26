package controller

import (
	"github.com/bitcapybara/cuckoo/core"
	"github.com/emirpasic/gods/sets/hashset"
	"sync"
)

type jobGroup struct {
	groups map[string]*hashset.Set
	mu     sync.Mutex
}

func newJobGroup() *jobGroup {
	return &jobGroup{
		groups: make(map[string]*hashset.Set),
	}
}

func (g *jobGroup) register(groupName string, addr core.NodeAddr) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if group, ok := g.groups[groupName]; !ok {
		g.groups[groupName] = hashset.New(addr)
	} else {
		group.Add(addr)
	}
}

func (g *jobGroup) unRegister(groupName string, addr core.NodeAddr) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if group, ok := g.groups[groupName]; ok {
		group.Remove(addr)
	}
}

func (g *jobGroup) getClients(groupName string) []core.NodeAddr {
	g.mu.Lock()
	defer g.mu.Unlock()
	values := g.groups[groupName].Values()
	result := make([]core.NodeAddr, len(values))
	for i, value := range values {
		result[i] = value.(core.NodeAddr)
	}
	return result
}
