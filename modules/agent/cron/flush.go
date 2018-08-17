package cron

import (
	"encoding/json"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/store"
	"time"
)

func updateStoreStatus() {
	_, err := funcs.GetData(g.Config().Backend.Check)
	if err == nil {
		store.UpdateStoreStatus(true)
	} else {
		store.UpdateStoreStatus(false)
	}
}

func consumeStore(queue chan string) {
	for {
		select {
		case data := <-queue:
			var mvs []*model.MetricValue
			if err := json.Unmarshal([]byte(data), &mvs); err == nil {
				g.SendToTransfer(mvs)
			}
		}
	}
}

func cleanStale() {
	for {
		s := store.GetStore()
		timestamp := time.Now().Unix() - int64(g.Config().Backend.Expire)
		s.CleanStale(timestamp)
		time.Sleep(time.Duration(g.Config().Backend.Expire) * time.Second)
	}
}

func eatStore(queue chan string) {
	for {
		s := store.GetStore()
		updateStoreStatus()
		if store.GetStoreStatus() {
			for data := s.Read(); data != ""; data = s.Read() {
				queue <- data
			}
		}
		time.Sleep(time.Second)
	}
}

func FlushStore() {
	if !g.Config().Backend.Enabled {
		return
	}
	queue := make(chan string, g.Config().Transfer.Interval)
	go eatStore(queue)
	go consumeStore(queue)
	go cleanStale()
}
