package router

import "github.com/bitcapybara/cuckoo/core"

type routeLast struct {
}

func newRouteLast() routeLast {
	return routeLast{}
}

func (r routeLast) route(groupName string, clients []core.NodeAddr) core.NodeAddr {
	return clients[len(clients)-1]
}
