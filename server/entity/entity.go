package entity

import (
	"github.com/bitcapybara/cuckoo/core"
	"time"
)

type JobInfo struct {
	Job    core.Job
	Next   time.Time
	Prev   time.Time
}

// CmdType 表示一个客户端请求类型
type CmdType uint8

const (
	Add    CmdType = iota // 添加任务
	Update                // 更新任务
	Delete                // 删除
)

type Cmd struct {
	CmdType CmdType
	Job core.Job
}
