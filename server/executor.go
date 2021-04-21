package main

// 任务的执行器，对应多个客户端
type Executor struct {
	name    string
	comment string
	clients []string
}
