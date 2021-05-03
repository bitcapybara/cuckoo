package router

import "github.com/bitcapybara/cuckoo/core"

type Router interface {
	route(string, []core.NodeAddr) core.NodeAddr
}

// 路由策略是有状态的，缓存在map中
var routers = map[core.RouteType]Router{
	core.First:  newRouteFirst(),
	core.Last:   newRouteLast(),
	core.Round:  newRouteRound(),
	core.Random: newRouteRandom(),
}

func Route(router core.RouteType, groupName string, clients []core.NodeAddr) core.NodeAddr {
	return routers[router].route(groupName, clients)
}
