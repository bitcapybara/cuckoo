package main

import (
	"fmt"
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/controller"
	"github.com/bitcapybara/cuckoo/server/entity"
	"github.com/bitcapybara/raft"
	"github.com/vmihailenco/msgpack/v5"
)

type Server struct {
	addr       string
	controller *controller.ScheduleController
}

func newServer(config raft.Config) *Server {
	return &Server{
		addr:       string(config.Peers[config.Me]),
		controller: controller.NewScheduleController(raft.NewNode(config)),
	}
}

func (s *Server) Start() {
	// 开启 raft 循环
	go s.controller.Start()
}

// 接收来自客户端的心跳注册请求
func (s *Server) Heartbeat(req core.HeartbeatReq, reply *core.CudReply) error {
	return nil
}

// 接收来自客户端的添加任务请求
func (s *Server) AddJob(req core.AddJobReq, reply *core.CudReply) error {
	cmd := entity.Cmd{
		CmdType: entity.Add,
		JobInfo: entity.JobInfo{
			Job:    req.Job,
			Enable: req.Enable,
		},
	}
	return s.sendApplyCommand(cmd, reply)
}

// 接收来自客户端的修改任务请求
func (s *Server) UpdateJob(req core.UpdateJobReq, reply *core.CudReply) error {
	cmd := entity.Cmd{
		CmdType: entity.Update,
		JobInfo: entity.JobInfo{
			Job: req.Job,
		},
	}
	return s.sendApplyCommand(cmd, reply)
}

// 接收来自客户端的删除任务请求
func (s *Server) DeleteJob(req core.DeleteJobReq, reply *core.CudReply) error {
	cmd := entity.Cmd{
		CmdType: entity.Delete,
		JobInfo: entity.JobInfo{
			Job: core.Job{Id: req.JobId},
		},
	}
	return s.sendApplyCommand(cmd, reply)
}

// 接收来自客户端的查看任务请求
func (s *Server) PageQuery(req core.PageQueryReq, reply *core.CudReply) error {
	return nil
}

func (s *Server) sendApplyCommand(cmd entity.Cmd, reply *core.CudReply) error {
	data, msErr := msgpack.Marshal(cmd)
	if msErr != nil {
		return fmt.Errorf("序列化失败：%w", msErr)
	}
	args := raft.ApplyCommand{
		Data: data,
	}
	var res raft.ApplyCommandReply
	raftErr := s.controller.Node.ApplyCommand(args, &res)
	if raftErr != nil {
		return fmt.Errorf("raft 操作失败！%w", raftErr)
	}
	reply.Status = core.Status(res.Status)
	reply.Leader = core.NodeAddr(res.Leader.Addr)
	return nil
}
