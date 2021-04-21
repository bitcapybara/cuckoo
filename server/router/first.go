package router

import "github.com/bitcapybara/cuckoo/core"

type routeFirst struct {
}

func newRouteFirst() routeFirst {
	return routeFirst{}
}

func (r routeFirst) route(clients []core.NodeAddr) core.NodeAddr {
	panic("implement me")
}


