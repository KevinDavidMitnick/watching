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
	ticker := time.NewTicker(time.Duration(g.Config().Transfer.Interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, err := funcs.GetData(g.Config().Backend.Check)
			if err == nil {
				store.UpdateStoreStatus(true)
			} else {
				store.UpdateStoreStatus(false)
			}
		}
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
	ticker := time.NewTicker(time.Duration(g.Config().Backend.Expire) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s := store.GetStore()
			timestamp := time.Now().Unix() - int64(g.Config().Backend.Expire)
			s.CleanStale(timestamp)
		}
	}
}

func eatStore(queue chan string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s := store.GetStore()
			for data := s.Read(); store.GetStoreStatus() && data != ""; data = s.Read() {
				queue <- data
			}
		}
	}
}

func FlushStore() {
	if !g.Config().Backend.Enabled {
		return
	}
	queue := make(chan string, g.Config().Transfer.Interval)
	go updateStoreStatus()
	go eatStore(queue)
	go consumeStore(queue)
	go cleanStale()
}
