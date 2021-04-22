package components

type scheduleFsm struct {
	executors map[string]Executor // 所有执行器，key=ExecutorName
	timeRing  TimeRing            // 时间轮
}
