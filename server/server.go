package main

import (
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/raft"
)

type Server struct {
	node      *raft.Node
	executors map[string]Executor
}

func NewServer() *Server {
	config := raft.Config{
	}
	node := raft.NewNode(config)
	return &Server{node: node}
}

// 接收来自客户端的心跳注册请求
func (c *Server) Heartbeat(req core.HeartbeatReq, reply *core.RpcReply) error {
	return nil
}

// 接收来自客户端的提交任务请求
func (c *Server) Submit(req core.SubmitReq, reply *core.RpcReply) error {
	return nil
}
