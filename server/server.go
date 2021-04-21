package main

import "github.com/bitcapybara/cuckoo/core"

type Server struct {

}

func NewServer() *Server {
	return &Server{}
}

// 接收来自客户端的心跳注册请求
func (c *Server) Heartbeat(req core.HeartbeatReq, reply *core.RpcReply) error {
	return nil
}

// 接收来自客户端的提交任务请求
func (c *Server) Submit(req core.SubmitReq, reply *core.RpcReply) error {
	return nil
}
