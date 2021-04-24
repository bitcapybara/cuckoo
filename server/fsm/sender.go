package fsm

import (
	"github.com/bitcapybara/cuckoo/core"
	"github.com/bitcapybara/cuckoo/server/entity"
)

// 网络通信接口
type JobSender interface {
	// 给客户端发送任务让其执行
	Send(clientAddr core.NodeAddr, info entity.JobInfo) error
}
