package controller

import (
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/raft"
	"time"
)

const (
	// 调度间隔时间
	ScheduleInterval = time.Second * 5
	// 每次调度最大任务数
	ScheduleMaxJob = 5000
)

type ScheduleController struct {
	Node      *raft.Node          // raft 节点
	executors map[string]Executor // 所有执行器，key=ExecutorName
	timeRing  *timeRing           // 时间轮，存放近期需要执行的任务，master节点使用
	jobPool   JobPool             // 任务池，存放所有任务
}

func NewScheduleController(node *raft.Node) *ScheduleController {
	return &ScheduleController{
		Node:     node,
		timeRing: NewTimeRing(),
	}
}

func (c *ScheduleController) Start() {
	// 开启 raft 循环
	go c.Node.Run()

	// 开启调度循环
	go c.runSchedule()

	// 开启时间轮循环
	go c.runTimeRing()
}

func (c *ScheduleController) runSchedule() {
	timer := time.NewTimer(ScheduleInterval)
	for c.Node.IsLeader() {
		func() {
			<-timer.C
			// 配置定时器
			now := time.Now()
			defer func() {
				end := time.Now()
				num := end.Sub(now).Milliseconds()%5000 + 1
				deviation := now.Add(ScheduleInterval * time.Duration(num)).Sub(end)
				timer.Reset(deviation)
			}()
			// 从任务池中获取未来 ScheduleInterval 时间内的 ScheduleMaxJob 个调度任务
			jobInfos := c.jobPool.Query(now.Add(ScheduleInterval), ScheduleMaxJob)
			if len(jobInfos) <= 0 {
				return
			}
			// 开始调度
			for _, jobInfo := range jobInfos {
				if now.After(jobInfo.Next) {
					// todo 错过了调度时间，立即执行一次
				}
				// 放入时间轮
				c.timeRing.put(jobInfo.Next.Second(), jobInfo.Job)
			}
			// 更新任务信息
			for _, jobInfo := range jobInfos {
				_ = c.jobPool.Update(jobInfo)
			}
		}()
	}
}

func (c *ScheduleController) runTimeRing() {
	timer := time.NewTimer(time.Second)
	for c.Node.IsLeader() {
		func() {
			<-timer.C
			// 配置定时器
			now := time.Now()
			defer func() {
				end := time.Now()
				num := end.Sub(now).Milliseconds()%5000 + 1
				deviation := now.Add(ScheduleInterval * time.Duration(num)).Sub(end)
				timer.Reset(deviation)
			}()
			// 取出时间轮上最近两秒的所有任务
			var ringItemData []core.Job
			for i := 0; i < 2; i++ {
				jobs := c.timeRing.getAndRemove((now.Second() + 60 - i) % 60)
				if len(jobs) > 0 {
					ringItemData = append(ringItemData, jobs...)
				}
			}
			if len(ringItemData) > 0 {
				// 触发任务
			}
		}()
	}
}
