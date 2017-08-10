package raft

type Ocean struct {
	islands	map[string]*Island
	stop chan int
}

func NewOcean() *Ocean {
	return &Ocean{
		islands: make(map[string]*Island),
		stop: make(chan int),
	}
}

func (o *Ocean) AddIsland(i *Island) {
	i.ocean = o
	o.islands[i.name] = i
}

func (o *Ocean) Stop() {
	o.stop <- 0
}

func (o *Ocean) Start() {
	for _, v := range o.islands {
		island := v
		go island.Run()
	}
	<-o.stop
}

