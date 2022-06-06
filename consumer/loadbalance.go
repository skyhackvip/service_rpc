package consumer

import (
	"math/rand"
	"time"
)

type LoadBalanceMode int

const (
	RandomBalance LoadBalanceMode = iota
	RoundRobinBalance
	WeightRoundRobinBalance
)

type LoadBalance interface {
	Get() string
}

func LoadBalanceFactory(mode LoadBalanceMode, servers []string) LoadBalance {
	switch mode {
	case RandomBalance:
		return newRandomBalance(servers)
	case RoundRobinBalance:
		return newRoundRobinBalance(servers)
	default:
		return newRandomBalance(servers)
	}
}

type randomBalance struct {
	servers []string
}

func newRandomBalance(servers []string) LoadBalance {
	return &randomBalance{servers: servers}
}

func (b *randomBalance) Get() string {
	rand.Seed(time.Now().Unix())
	return b.servers[rand.Intn(len(b.servers))]
}

type roundRobinBalance struct {
	servers []string
	curIdx  int
}

func newRoundRobinBalance(servers []string) LoadBalance {
	return &roundRobinBalance{servers: servers, curIdx: 0}
}

func (b *roundRobinBalance) Get() string {
	lens := len(b.servers)
	if b.curIdx >= lens {
		b.curIdx = 0
	}
	server := b.servers[b.curIdx]
	b.curIdx = (b.curIdx + 1) % lens
	return server
}
