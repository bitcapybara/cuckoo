package router

import "github.com/bitcapybara/cuckoo/core"

type routeFirst struct {
}

func newRouteFirst() routeFirst {
	return routeFirst{}
}

func (r routeFirst) route(groupName string, clients []core.NodeAddr) core.NodeAddr {
	return clients[0]
}


