package router

import "github.com/bitcapybara/cuckoo/core"

type Router interface {
	route([]core.NodeAddr) core.NodeAddr
}

var routers = map[core.RouteType]Router{
	core.First:  newRouteFirst(),
	core.Last:   newRouteLast(),
	core.Round:  newRouteRound(),
	core.Random: newRouteRandom(),
}

func GetRouter(router core.RouteType) Router {
	return routers[router]
}
