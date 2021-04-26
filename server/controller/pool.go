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

// 实现对象的序列化功能
type Serializable interface {

	Encode() ([]byte, error)

	Decode([]byte) error
}

// JobPool 是存放所有任务的任务池
// 每个节点维护一份作为多个备份
type JobPool interface {

	Serializable

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

// SliceJobPool 是以数组形式实现的 JobPool
type SliceJobPool struct {
	logger   raft.Logger
	listData *arraylist.List
	mu       sync.Mutex
}

func (s *SliceJobPool) Encode() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return msgpack.Marshal(s.listData)
}

func (s *SliceJobPool) Decode(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return msgpack.Unmarshal(data, s.listData)
}

func jobInfoComparator(a, b interface{}) int {
	info1 := a.(*entity.JobInfo)
	info2 := b.(*entity.JobInfo)
	return info1.Next.Second() - info2.Next.Second()
}

func NewSliceJobStorage(logger raft.Logger) *SliceJobPool {
	return &SliceJobPool{
		logger:   logger,
		listData: arraylist.New(),
	}
}

// 实现 JobPool 接口

func (s *SliceJobPool) Add(jobInfo *entity.JobInfo) error {
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

func (s *SliceJobPool) GetById(id core.JobId) *entity.JobInfo {
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

func (s *SliceJobPool) Update(jobInfo *entity.JobInfo) error {
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

func (s *SliceJobPool) DeleteById(jobId core.JobId) error {
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

func (s *SliceJobPool) Query(timeBefore time.Time, count int) []*entity.JobInfo {
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
