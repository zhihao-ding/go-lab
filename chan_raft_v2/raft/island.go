package raft

import (
	"fmt"
	"time"
	"math/rand"
)

const (
	ELECTION_TIMEOUT_TOP = 5000
	ELECTION_TIMEOUT_BOTTOM = 1000
	HEARTBEAT_TIMEOUT = 1000
)

type Island struct {
	name	string
	inPort	chan *Raft
	state	StateType
	ocean	*Ocean
}

func NewIsland(islandName string) *Island {
	return &Island{
		name: islandName,
		state: Candidate,
		inPort: make(chan *Raft),
	}
}

func (i *Island) info(str string) {
	fmt.Printf("%10s@%9s> %s\n", i.name, i.state, str)
}

func (i *Island) Run() {
	for {
		switch i.state {
		case Candidate:
			i.doCandidate()
		case Follower:
			i.doFollower()
		case Leader:
			i.doLeader()
		default:
			fmt.Printf("%10s at Unknown state!\n", i.name)
			time.Sleep(time.Second * 10)
		}
	}
}

func (i *Island) InPort(raft *Raft) {
	defer func() {
		if err := recover(); err != nil {
			i.info("InPort is closed, drop raft!")
		}
	}()
	i.inPort <- raft
}

func (i *Island) SendRaft(rt RaftType, to string) {
	i.ocean.islands[to].InPort(&Raft{id: rt, src: i.name})
}

func (i *Island) SpreadRaft(rt RaftType) int {
	count := 0
	for _, v := range i.ocean.islands {
		island := v
		if island.name != i.name {
			island.InPort(&Raft{id: rt, src: i.name})
			count += 1
		}
	}
	return count
}

func (i *Island) doCandidate() {
	stop := make(chan int)

	go func() {
		defer func() {
			stop <- 0
		}()

		for {
			select {
			case raft := <-i.inPort:
				i.info(fmt.Sprintf("Receive %s from %s.", raft.id, raft.src))
				switch raft.id {
				case Vote:
					i.SendRaft(VoteAck, raft.src)
					i.info("Change state to Follower.")
					i.state = Follower
					return
				case VoteAck:
					i.info("Change state to Leader.")
					i.state = Leader
					return
				default:
				}
			default:
			}
		}
	}()

	//rand.Seed(time.Now().Unix())
	randTimeout := rand.Intn(ELECTION_TIMEOUT_TOP) // - ELECTION_TIMEOUT_BOTTOM) + ELECTION_TIMEOUT_BOTTOM
	i.info(fmt.Sprintf("Rand election timeout: %d", randTimeout))
	t := time.NewTimer(time.Duration(randTimeout) * time.Millisecond)

	for {
		select {
		case <-t.C:
			fmt.Println()
			i.info("Election timeout, request votes...")
			i.SpreadRaft(Vote)
		case <-stop:
			//i.info("Receive out.")
			t.Stop()
			return
		default:
		}
	}
}

func (i *Island) doFollower() {
	for {
		t := time.NewTimer(time.Second * 2)
		select {
		case <-t.C:
			i.info("Waiting HeartBeat Timeout, Change state to Candidate.")
			i.state = Candidate
		    t.Stop()
			return
		case raft := <-i.inPort:
			i.info(fmt.Sprintf("Receive %s from %s.", raft.id, raft.src))
			if raft.id == HeartBeat {
				i.SendRaft(HeartBeatAck, raft.src)
			}
		}
		t.Stop()
	}
}

func (i *Island) doLeader() {
	dead := make(chan int)

	go func() {
		for {
			select {
			case <-i.inPort:
				//i.info(fmt.Sprintf("Receive %s from %s.", raft.id, raft.src))
			default:
			}
		}
	}()

	count := 0
	for {
		time.Sleep(time.Second * 1)
		fmt.Println()
		i.info("Send HeartBeat Message to others.")
		i.SpreadRaft(HeartBeat)
		count += 1
		if count > 10 {
			break
		}
	}

	i.info("Dead")
	<-dead
}

