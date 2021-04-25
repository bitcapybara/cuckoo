package core

type NodeAddr string

type HeartbeatReq struct {
	ExecutorName string
	LocalAddr    NodeAddr
}

type AddJobReq struct {
	Job
	Enable bool
}

type UpdateJobReq struct {
	Job
}

type DeleteJobReq struct {
	JobId
}

type PageQueryReq struct {
	JobId
	ExecutorName
	PageNum  int
	PageSize int
}

type Status uint8

type CudReply struct {
	Status   Status
	Leader   NodeAddr
}

type QueryReply struct {
	Jobs []Job
}
