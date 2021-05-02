package controller

import (
	"github.com/bitcapybara/cuckoo/core"
	"log"
	"sync"
	"time"
)

type jobGroup struct {
	groups map[string]map[core.NodeAddr]time.Time
	mu     sync.Mutex

	expire time.Duration
}

func newJobGroup(aliveExpire time.Duration) *jobGroup {
	return &jobGroup{
		groups: make(map[string]map[core.NodeAddr]time.Time),
		expire: aliveExpire,
	}
}

func (g *jobGroup) register(groupName string, addr core.NodeAddr) {
	g.mu.Lock()
	defer g.mu.Unlock()
	log.Println("节点注册，group="+groupName+","+"addr="+string(addr))
	if group, ok := g.groups[groupName]; !ok {
		g.groups[groupName] = map[core.NodeAddr]time.Time{addr: time.Now().Add(g.expire)}
	} else {
		group[addr] = time.Now().Add(g.expire)
	}
}

func (g *jobGroup) unRegister(groupName string, addr core.NodeAddr) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if group, gOk := g.groups[groupName]; gOk {
		if _, aOk := group[addr]; aOk {
			delete(group, addr)
		}
	}
}

func (g *jobGroup) getClients(groupName string) []core.NodeAddr {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := time.Now()
	values := g.groups[groupName]
	result := make([]core.NodeAddr, 0)
	for addr, ex := range values {
		if ex.Before(now) {
			delete(values, addr)
			continue
		}
		result = append(result, addr)
	}
	return result
}
