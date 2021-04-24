package fsm

type scheduleFsm struct {
	executors map[string]Executor // 所有执行器，key=ExecutorName
	timeRing  TimeRing            // 时间轮，存放近期需要执行的任务

}
