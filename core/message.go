package core

type NodeAddr string

type HeartbeatReq struct {
	Group     string
	LocalAddr NodeAddr
}

type HeartbeatReply struct {
	Status Status
	Leader NodeAddr
}

type AddJobReq struct {
	Job
}

type UpdateJobReq struct {
	Job
}

type DeleteJobReq struct {
	JobId
}

type PageQueryReq struct {
	JobId
	Group    string
	PageNum  int
	PageSize int
}

type Status uint8

const (
	NotLeader Status = iota
	Ok
)

type CudReply struct {
	Status Status
	Leader NodeAddr
}

type QueryReply struct {
	Jobs []Job
}
