package controller

import (
	"fmt"
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/router"
	"github.com/bitcapybara/cuckoo/server/schedule"
	"github.com/bitcapybara/raft"
	"sync"
	"time"
)

const (
	// 调度间隔时间
	ScheduleInterval = time.Second * 5
	// 每次调度最大任务数
	ScheduleMaxJob = 5000
)

// 网络通信接口
type JobDispatcher interface {
	// 给客户端发送任务让其执行
	Dispatch(clientAddr core.NodeAddr, jobId core.JobId) error
}

type Config struct {
	Node           *raft.Node
	JobPool        JobPool
	Logger         raft.Logger
	Dispatcher     JobDispatcher
	ExecutorExpire time.Duration
}

type ScheduleController struct {
	logger      raft.Logger   // 日志
	Node        *raft.Node    // raft 节点
	jobGroup    *jobGroup     // 执行器对应的客户端节点列表 // 所有执行器，key=JobGroupName
	timeRing    *timeRing     // 时间轮，存放近期需要执行的任务，master节点使用
	jobPool     JobPool       // 任务池，存放所有任务
	dispatcher  JobDispatcher // 任务分发器
	initialized bool          // 是否已初始化
	mu          sync.Mutex
}

func NewScheduleController(config Config) *ScheduleController {
	return &ScheduleController{
		logger:     config.Logger,
		Node:       config.Node,
		jobGroup:   newJobGroup(config.ExecutorExpire),
		timeRing:   NewTimeRing(),
		jobPool:    config.JobPool,
		dispatcher: config.Dispatcher,
	}
}

func (s *ScheduleController) Start() {
	// 开启 raft 循环
	go s.Node.Run()

	// 注册角色变更观察者
	roleObserver := make(chan raft.RoleStage)
	s.Node.AddRoleObserver(roleObserver)

	schedTimer := time.NewTimer(ScheduleInterval)
	ringTimer := time.NewTimer(time.Second)
	for {
		if s.Node.IsLeader() {
			if !s.initialized {
				s.init()
			}
			select {
			case <-schedTimer.C:
				go s.runSchedule(schedTimer)
			case <-ringTimer.C:
				go s.runTimeRing(ringTimer)
			}
		} else {
			role := <-roleObserver
			if role != raft.Leader {
				s.mu.Lock()
				s.initialized = false
				s.mu.Unlock()
			}
		}
	}
}

func (s *ScheduleController) init() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	option := QueryOption{
		timeBefore: now,
	}
	infos := s.jobPool.Query(option)
	for _, info := range infos {
		if info.Next.IsZero() {
			info.Next = schedule.Schedule(info.Job.ScheduleRule, now)
		} else {
			info.Next = schedule.Schedule(info.Job.ScheduleRule, info.Next)
		}
		_ = s.jobPool.Update(info)
	}
	s.initialized = true
}

func (s *ScheduleController) runSchedule(timer *time.Timer) {
	// 配置定时器
	now := time.Now()
	defer func() {
		end := time.Now()
		num := end.Sub(now).Milliseconds()%5000 + 1
		deviation := now.Add(ScheduleInterval * time.Duration(num)).Sub(end)
		timer.Reset(deviation)
	}()
	// 从任务池中获取未来 ScheduleInterval 时间内的 ScheduleMaxJob 个调度任务
	option := QueryOption{
		timeBefore: now.Add(ScheduleInterval),
		count:      ScheduleMaxJob,
	}
	jobInfos := s.jobPool.Query(option)
	if len(jobInfos) <= 0 {
		return
	}
	// 开始调度
	for _, jobInfo := range jobInfos {
		if now.After(jobInfo.Next) {
			// 错过了调度时间，立即执行一次
			s.Trigger(jobInfo.Job)
		}
		// 放入时间轮
		s.timeRing.put(jobInfo.Next.Second(), jobInfo.Job)
		jobInfo.Next = schedule.Schedule(jobInfo.Job.ScheduleRule, jobInfo.Next)
	}
	// 更新任务信息
	for _, jobInfo := range jobInfos {
		_ = s.jobPool.Update(jobInfo)
	}
}

func (s *ScheduleController) runTimeRing(timer *time.Timer) {
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
		jobs := s.timeRing.getAndRemove((now.Second() + 60 - i) % 60)
		if len(jobs) > 0 {
			ringItemData = append(ringItemData, jobs...)
		}
	}
	if len(ringItemData) > 0 {
		// 触发任务
		for _, job := range ringItemData {
			s.Trigger(job)
		}
	}
}

func (s *ScheduleController) Trigger(job core.Job) {
	go func() {
		clientAddr := router.Route(job.Router, s.jobGroup.getClients(job.Group))
		dispatchErr := s.dispatcher.Dispatch(clientAddr, job.Id)
		if dispatchErr != nil {
			s.logger.Error(fmt.Errorf("分发任务出错：%w", dispatchErr).Error())
		}
	}()
}

func (s *ScheduleController) Register(groupName string, addr core.NodeAddr) {
	s.jobGroup.register(groupName, addr)
}
