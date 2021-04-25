package client

import (
	"github.com/bitcapybara/cuckoo/core"
)

type RpcClient interface {
	Heartbeat(addr core.NodeAddr, heartbeatReq core.HeartbeatReq, reply *core.CudReply) error
	Submit(addr core.NodeAddr, submitReq core.AddJobReq, reply *core.CudReply) error
}
