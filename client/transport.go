package client

import (
	"github.com/bitcapybara/cuckoo/core"
)

type Transport interface {
	Heartbeat(addr core.NodeAddr, heartbeatReq core.HeartbeatReq, reply *core.RpcReply) error
	Submit(addr core.NodeAddr, submitReq core.SubmitReq, reply *core.RpcReply) error
}
