package main

import (
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/raft"
)

type Server struct {
	addr string
	node *raft.Node
}

func newServer(role raft.RoleStage, me raft.NodeId, peers map[raft.NodeId]raft.NodeAddr) *Server {
	config := raft.Config{
		Peers:              peers,
		Me:                 me,
		Role:               role,
		ElectionMaxTimeout: 10000,
		ElectionMinTimeout: 5000,
		HeartbeatTimeout:   1000,
		MaxLogLength:       50,
	}

	return &Server{
		node: raft.NewNode(config),
	}
}

func (s *Server) Start() {
	// 开启 raft 循环
	go s.node.Run()
}

// 接收来自客户端的心跳注册请求
func (c *Server) Heartbeat(req core.HeartbeatReq, reply *core.RpcReply) error {
	return nil
}

// 接收来自客户端的提交任务请求
func (c *Server) Submit(req core.SubmitReq, reply *core.RpcReply) error {
	return nil
}
