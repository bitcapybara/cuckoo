package main

import (
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/controller"
	"github.com/bitcapybara/raft"
)

type Server struct {
	addr       string
	controller *controller.ScheduleController
}

func newServer(config raft.Config) *Server {
	return &Server{
		addr: string(config.Peers[config.Me]),
		controller: controller.NewScheduleController(raft.NewNode(config)),
	}
}

func (s *Server) Start() {
	// 开启 raft 循环
	go s.controller.Start()
}

// 接收来自客户端的心跳注册请求
func (s *Server) Heartbeat(req core.HeartbeatReq, reply *core.RpcReply) error {
	return nil
}

// 接收来自客户端的添加任务请求
func (s *Server) AddJob(req core.SubmitReq, reply *core.RpcReply) error {
	return nil
}

// 接收来自客户端的修改任务请求
func (s *Server) UpdateJob() error {
	return nil
}

// 接收来自客户端的删除任务请求
func (s *Server) DeleteJob() error {
	return nil
}

// 接收来自客户端的查看任务请求
func (s *Server) PageList() error {
	return nil
}
