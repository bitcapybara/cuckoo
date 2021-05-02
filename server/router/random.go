package router

import (
	"github.com/bitcapybara/cuckoo/core"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type routeRandom struct {
}

func newRouteRandom() routeRandom {
	return routeRandom{}
}

func (r routeRandom) route(groupName string, clients []core.NodeAddr) core.NodeAddr {
	return clients[rand.Intn(len(clients))]
}

