package main

import (
	"chan_raft_v2/raft"
)

func main() {
	o := raft.NewOcean()
	o.AddIsland(raft.NewIsland("China"))
	o.AddIsland(raft.NewIsland("America"))
	o.AddIsland(raft.NewIsland("England"))
	o.Start()
}
