package fsm

import (
	"errors"
	"fmt"
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/entity"
	"github.com/emirpasic/gods/lists/arraylist"
	"time"
)

// JobPool 是存放所有任务的任务池
type JobPool interface {
	// 添加任务到持久化存储
	Add(jobInfo *entity.JobInfo) error

	// 获取某一id的任务
	GetById(id core.JobId) *entity.JobInfo

	// 更新任务信息
	Update(jobInfo *entity.JobInfo) error

	// 删除某个任务
	DeleteById(id core.JobId) error

	// 获取要调度的任务
	JobQuery(timeBefore time.Time, count int) []*entity.JobInfo
}

// 以数组形式实现的 JobPool
type SliceJobStorage struct {
	listData *arraylist.List
}

func jobInfoComparator(a, b interface{}) int {
	info1 := a.(*entity.JobInfo)
	info2 := b.(*entity.JobInfo)
	return info1.Next.Second() - info2.Next.Second()
}

func NewSliceJobStorage() *SliceJobStorage {
	return &SliceJobStorage{
		listData: arraylist.New(),
	}
}

func (s *SliceJobStorage) Add(jobInfo *entity.JobInfo) error {
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

func (s *SliceJobStorage) JobQuery(timeBefore time.Time, count int) []*entity.JobInfo {
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
