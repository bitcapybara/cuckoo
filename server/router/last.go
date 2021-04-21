package router

import "github.com/bitcapybara/cuckoo/core"

type routeLast struct {
}

func newRouteLast() routeLast {
	return routeLast{}
}

func (r routeLast) route(clients []core.NodeAddr) core.NodeAddr {
	panic("implement me")
}
