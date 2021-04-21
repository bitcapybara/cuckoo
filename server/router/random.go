package router

import "github.com/bitcapybara/cuckoo/core"

type routeRandom struct {
}

func newRouteRandom() routeRandom {
	return routeRandom{}
}

func (r routeRandom) route(clients []core.NodeAddr) core.NodeAddr {
	panic("implement me")
}

