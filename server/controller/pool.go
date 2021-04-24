package controller

import (
	"errors"
	"fmt"
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/entity"
	"github.com/bitcapybara/raft"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/vmihailenco/msgpack/v5"
	"sync"
	"time"
)

// JobPool 是存放所有任务的任务池
// 每个节点维护一份作为多个备份
type JobPool interface {
	// Add 添加任务到持久化存储
	Add(jobInfo *entity.JobInfo) error

	// GetById 获取某一id的任务
	GetById(id core.JobId) *entity.JobInfo

	// Update 更新任务信息
	Update(jobInfo *entity.JobInfo) error

	// DeleteById 删除某个任务
	DeleteById(id core.JobId) error

	// Query 获取要调度的任务
	Query(timeBefore time.Time, count int) []*entity.JobInfo
}

// SliceJobStorage 是以数组形式实现的 JobPool
type SliceJobStorage struct {
	logger   raft.Logger
	listData *arraylist.List
	mu       sync.Mutex
}

func jobInfoComparator(a, b interface{}) int {
	info1 := a.(*entity.JobInfo)
	info2 := b.(*entity.JobInfo)
	return info1.Next.Second() - info2.Next.Second()
}

func NewSliceJobStorage(logger raft.Logger) *SliceJobStorage {
	return &SliceJobStorage{
		logger:   logger,
		listData: arraylist.New(),
	}
}

// 实现 JobPool 接口

func (s *SliceJobStorage) Add(jobInfo *entity.JobInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if jobInfo.Next.IsZero() {
		return errors.New("未指定任务的下次执行时间")
	}
	if info := s.GetById(jobInfo.Job.Id); info != nil {
		return errors.New("当前任务已存在")
	}
	s.listData.Add(jobInfo)
	s.listData.Sort(jobInfoComparator)
	return nil
}

func (s *SliceJobStorage) GetById(id core.JobId) *entity.JobInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	index, v := s.listData.Find(func(index int, value interface{}) bool {
		info := value.(*entity.JobInfo)
		return info.Job.Id == id
	})
	if index == -1 {
		return nil
	}
	return v.(*entity.JobInfo)
}

func (s *SliceJobStorage) Update(jobInfo *entity.JobInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if jobInfo.Next.IsZero() {
		return errors.New("未指定任务的下次执行时间")
	}
	jobId := jobInfo.Job.Id
	index, _ := s.listData.Find(func(index int, value interface{}) bool {
		info := value.(*entity.JobInfo)
		return info.Job.Id == jobId
	})
	if index != -1 {
		s.listData.Set(index, jobInfo)
	} else {
		return errors.New(fmt.Sprintf("没有找到 id=%d 的任务", jobId))
	}
	s.listData.Sort(jobInfoComparator)
	return nil
}

func (s *SliceJobStorage) DeleteById(jobId core.JobId) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	index, _ := s.listData.Find(func(index int, value interface{}) bool {
		info := value.(*entity.JobInfo)
		return info.Job.Id == jobId
	})
	if index != -1 {
		s.listData.Remove(index)
	}
	s.listData.Sort(jobInfoComparator)
	return nil
}

func (s *SliceJobStorage) Query(timeBefore time.Time, count int) []*entity.JobInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.listData.Select(func(index int, value interface{}) bool {
		info := value.(*entity.JobInfo)
		return info.Next.Before(timeBefore) && index < count
	})
	result := make([]*entity.JobInfo, list.Size())
	iterator := list.Iterator()
	for i := 0; iterator.Next(); i++ {
		result[i] = iterator.Value().(*entity.JobInfo)
	}
	return result
}

// 实现 raft.Fsm 接口
// 结构体的编码/解码使用 msgPack

func (s *SliceJobStorage) Apply(data []byte) error {
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
		return s.Add(jobInfo)
	case entity.Update:
		return s.Update(jobInfo)
	case entity.Delete:
		return s.DeleteById(jobInfo.Job.Id)
	default:
		err := errors.New(fmt.Sprintf("不支持的命令：%d", cmd.CmdType))
		s.logger.Error(err.Error())
		return err
	}
}

func (s *SliceJobStorage) Serialize() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, umsErr := msgpack.Marshal(s.listData)
	if umsErr != nil {
		err := fmt.Errorf("序列化任务池失败！%w", umsErr)
		s.logger.Error(err.Error())
		return nil, err
	}
	return data, nil
}

func (s *SliceJobStorage) Install(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list arraylist.List
	msErr := msgpack.Unmarshal(data, &list)
	if msErr != nil {
		err := fmt.Errorf("反序列化状态机命令失败！%w", msErr)
		s.logger.Error(err.Error())
		return err
	}
	s.listData = &list
	return nil
}
