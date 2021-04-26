package controller

import (
	"errors"
	"fmt"
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/entity"
	"github.com/bitcapybara/raft"
	"github.com/vmihailenco/msgpack/v5"
	"sync"
	"time"
)

const (
	// 调度间隔时间
	ScheduleInterval = time.Second * 5
	// 每次调度最大任务数
	ScheduleMaxJob = 5000
)

type ScheduleController struct {
	logger    raft.Logger
	Node      *raft.Node          // raft 节点
	executors map[string]Executor // 所有执行器，key=ExecutorName
	timeRing  *timeRing           // 时间轮，存放近期需要执行的任务，master节点使用
	jobPool   JobPool             // 任务池，存放所有任务
	mu        sync.Mutex
}

func NewScheduleController(node *raft.Node, jobPool JobPool, logger raft.Logger) *ScheduleController {
	return &ScheduleController{
		logger: logger,
		Node:     node,
		timeRing: NewTimeRing(),
		jobPool: jobPool,
	}
}

func (s *ScheduleController) Start() {
	// 开启 raft 循环
	go s.Node.Run()

	// 开启调度循环
	go s.runSchedule()

	// 开启时间轮循环
	go s.runTimeRing()
}

func (s *ScheduleController) runSchedule() {
	timer := time.NewTimer(ScheduleInterval)
	for s.Node.IsLeader() {
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
			jobInfos := s.jobPool.Query(now.Add(ScheduleInterval), ScheduleMaxJob)
			if len(jobInfos) <= 0 {
				return
			}
			// 开始调度
			for _, jobInfo := range jobInfos {
				if now.After(jobInfo.Next) {
					// todo 错过了调度时间，立即执行一次
				}
				// 放入时间轮
				s.timeRing.put(jobInfo.Next.Second(), jobInfo.Job)
			}
			// 更新任务信息
			for _, jobInfo := range jobInfos {
				_ = s.jobPool.Update(jobInfo)
			}
		}()
	}
}

func (s *ScheduleController) runTimeRing() {
	timer := time.NewTimer(time.Second)
	for s.Node.IsLeader() {
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
				jobs := s.timeRing.getAndRemove((now.Second() + 60 - i) % 60)
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

// 实现 raft.Fsm 接口
// 结构体的编码/解码使用 msgPack

func (s *ScheduleController) Apply(data []byte) error {
	var cmd entity.Cmd
	msErr := msgpack.Unmarshal(data, &cmd)
	if msErr != nil {
		err := fmt.Errorf("反序列化状态机命令失败！%w", msErr)
		s.logger.Error(err.Error())
		return err
	}
	jobInfo := &cmd.JobInfo
	switch cmd.CmdType {
	case entity.Add:
		return s.jobPool.Add(jobInfo)
	case entity.Update:
		return s.jobPool.Update(jobInfo)
	case entity.Delete:
		return s.jobPool.DeleteById(jobInfo.Job.Id)
	default:
		err := errors.New(fmt.Sprintf("不支持的命令：%d", cmd.CmdType))
		s.logger.Error(err.Error())
		return err
	}
}

func (s *ScheduleController) Serialize() ([]byte, error) {
	data, umsErr := s.jobPool.Encode()
	if umsErr != nil {
		err := fmt.Errorf("序列化任务池失败！%w", umsErr)
		s.logger.Error(err.Error())
		return nil, err
	}
	return data, nil
}

func (s *ScheduleController) Install(data []byte) error {
	msErr := s.jobPool.Decode(data)
	if msErr != nil {
		err := fmt.Errorf("反序列化状态机命令失败！%w", msErr)
		s.logger.Error(err.Error())
		return err
	}
	return nil
}
