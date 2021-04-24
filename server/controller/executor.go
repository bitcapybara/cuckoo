package controller

// 任务的执行器
type Executor struct {
	name    string  // 执行器名称
	comment string  // 备注
	clients []string // 执行器对应的客户端节点列表
}
