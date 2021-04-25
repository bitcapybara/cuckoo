package client

import "github.com/bitcapybara/cuckoo/core"

type CuckooClient struct {
	transport RpcClient        // 网络通信接口，客户端调用此接口发送网络请求
	local     core.NodeAddr    // 客户端本地地址
	remote    core.NodeAddr // 服务端所有节点地址
}

// 创建客户端
func NewCuckooClient() *CuckooClient {
	client := &CuckooClient{}
	client.start()
	return client
}

// 客户端循环向服务器发送心跳
func (c *CuckooClient) start() {
	go func() {
		for {
			_ = c.transport.Heartbeat(c.remote, core.HeartbeatReq{}, &core.CudReply{})
		}
	}()
}

func (c *CuckooClient) Submit(job core.Job) {
	_ = c.transport.Submit(c.remote, core.AddJobReq{}, &core.CudReply{})
}
