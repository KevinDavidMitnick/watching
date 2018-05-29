package g

import (
	"sync"
	"time"
)

type SafeStrategyGroupMap struct {
	sync.RWMutex
	M map[int][]int
}

var (
	HbsClient        *SingleConnRpcClient
	StrategyGroupMap = &SafeStrategyGroupMap{M: make(map[int][]int)}
)

func InitHbsClient() {
	HbsClient = &SingleConnRpcClient{
		RpcServers: Config().Hbs.Servers,
		Timeout:    time.Duration(Config().Hbs.Timeout) * time.Millisecond,
	}
}

func (this *SafeStrategyGroupMap) ReInit(m map[int][]int) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeStrategyGroupMap) Get() map[int][]int {
	this.RLock()
	defer this.RUnlock()
	return this.M
}
