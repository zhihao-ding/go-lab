package raft

type StateType int

const (
	Candidate = iota
	Follower
	Leader
)

var stateTypes = [...]string {
	"Candidate",
	"Follower",
	"Leader",
}

func (t StateType) String() string {
	return stateTypes[t]
}
