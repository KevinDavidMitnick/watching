package cron

import (
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/store"
	"time"
)

func UpdateStoreStatus() {
	_, err := funcs.GetData(g.Config().Backend.Check)
	if err == nil {
		store.UpdateStoreStatus(true)
	} else {
		store.UpdateStoreStatus(false)
	}
}

func FlushStore() {
	if !g.Config().Backend.Enabled {
		return
	}
	interval := time.Duration(g.Config().Transfer.Interval)
	s := store.GetStore()
	defer s.Close()
	for {
		timestamp := time.Now().Unix() - int64(g.Config().Backend.Expire)
		data := make([]map[string]interface{}, 0)
		queue := make(chan string, g.Config().Transfer.Interval)
		s.CleanStale(timestamp, data)
		s = store.GetStore()
		UpdateStoreStatus()

		go func() {
			for {
				select {
				case data := <-queue:
					if data == "" {
						break
					}
					funcs.SubmitData(g.Config().Backend.Backup, []byte(data))
				}
			}
		}()
		for data := s.Read(); store.GetStoreStatus() && data != ""; data = s.Read() {
			queue <- data
		}
		close(queue)
		time.Sleep(interval * time.Second)
	}
}
