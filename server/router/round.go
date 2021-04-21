package router

import "github.com/bitcapybara/cuckoo/core"

type routeRound struct {
}

func newRouteRound() routeRound {
	return routeRound{}
}

func (r routeRound) route(clients []core.NodeAddr) core.NodeAddr {
	panic("implement me")
}
