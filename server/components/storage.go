package components

import (
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/entity"
	"time"
)

// 任务持久化存储
type JobStorage interface {
	// 添加任务到持久化存储
	Add(jobInfo entity.JobInfo) error

	// 获取某一id的任务
	GetById(id core.JobId) entity.JobInfo

	// 更新任务信息
	Update(jobInfo entity.JobInfo) error

	// 删除某个任务
	DeleteById(id core.JobId) error

	// 获取要调度的任务
	JobQuery(timeBefore time.Time, count int) []entity.JobInfo
}
