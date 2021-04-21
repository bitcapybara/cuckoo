package client

import "github.com/bitcapybara/cuckoo/core"

type CuckooClient struct {
	transport Transport
	local     core.NodeAddr
	remote    *core.RemoteInfo
}

func NewCuckooClient() *CuckooClient {
	return &CuckooClient{}
}

func (c *CuckooClient) heartbeat() {
	_ = c.transport.Heartbeat(c.remote.LeaderAddr, core.HeartbeatReq{}, &core.RpcReply{})
}

func (c *CuckooClient) Submit(job core.Job) {
	_ = c.transport.Submit(c.remote.LeaderAddr, core.SubmitReq{}, &core.RpcReply{})
}
