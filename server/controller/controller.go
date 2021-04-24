package controller

import "github.com/bitcapybara/raft"

type ScheduleController struct {
	node      *raft.Node          // raft 节点
	executors map[string]Executor // 所有执行器，key=ExecutorName
	timeRing  TimeRing            // 时间轮，存放近期需要执行的任务，master节点使用
	jobPool   JobPool             // 任务池，存放所有任务
}

func NewScheduleController(node *raft.Node) *ScheduleController {
	return &ScheduleController{node: node}
}

func (c *ScheduleController) Start() {

}
