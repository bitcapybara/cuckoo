package server

import (
	"fmt"
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/controller"
	"github.com/bitcapybara/cuckoo/server/entity"
	"github.com/bitcapybara/raft"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type Config struct {
	RaftConfig     raft.Config
	JobPool        controller.JobPool
	JobDispatcher  controller.JobDispatcher
	ExecutorExpire time.Duration
}

type Server struct {
	addr       string
	controller *controller.ScheduleController
}

func NewServer(config Config) *Server {
	raftConfig := config.RaftConfig
	raftNode := raft.NewNode(raftConfig)
	var jobPool controller.JobPool
	if config.JobPool != nil {
		jobPool = config.JobPool
	} else {
		jobPool = controller.NewSliceJobPool(raftConfig.Logger)
	}
	raftConfig.Fsm = controller.NewJobPoolFsm(raftConfig.Logger, jobPool)
	cc := controller.Config{
		Node: raftNode,
		JobPool: jobPool,
		Logger: raftConfig.Logger,
		Dispatcher: config.JobDispatcher,
		ExecutorExpire: config.ExecutorExpire,
	}
	return &Server{
		addr:       string(raftConfig.Peers[raftConfig.Me]),
		controller: controller.NewScheduleController(cc),
	}
}

func (s *Server) Start() {
	// 开启主循环
	s.controller.Start()
}

// 接收来自客户端的心跳注册请求
func (s *Server) Heartbeat(req core.HeartbeatReq, reply *core.CudReply) error {
	if !s.controller.Node.IsLeader() {
		reply.Status = core.NotLeader
		reply.Leader = core.NodeAddr(s.controller.Node.GetLeader())
	} else {
		s.controller.Register(req.Group, req.LocalAddr)
	}
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
