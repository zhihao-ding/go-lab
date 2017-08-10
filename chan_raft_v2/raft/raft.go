package raft

type RaftType int

const (
	Vote = iota
	VoteAck
	HeartBeat
	HeartBeatAck
)

var raftTypes = [...]string {
	"Vote",
	"VoteAck",
	"HeartBeat",
	"HeartBeatAck",
}

func (t RaftType) String() string {
	return raftTypes[t]
}

type Raft struct {
	id	RaftType
	src	string
}

