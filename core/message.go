package core

type NodeAddr string

type HeartbeatReq struct {
	ExecutorName string
	LocalAddr    NodeAddr
}

type SubmitReq struct {

}

type Status uint8

const (
	Ok Status = iota
	NotLeader
)

type RpcReply struct {
	Status Status
	Remote RemoteInfo
}
