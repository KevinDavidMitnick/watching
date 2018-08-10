// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"time"

	"encoding/json"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/funcs"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/store"
)

func InitDataHistory() {
	for {
		funcs.UpdateCpuStat()
		funcs.UpdateDiskStats()
		time.Sleep(g.COLLECT_INTERVAL)
	}
}

func Collect() {

	if !g.Config().Transfer.Enabled {
		return
	}

	if len(g.Config().Transfer.Addrs) == 0 {
		return
	}

	for _, v := range funcs.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*model.MetricValue) {
	t := time.NewTicker(time.Second * time.Duration(sec))
	defer t.Stop()
	for {
		<-t.C

		for _, fn := range fns {
			go collectMetric(sec, fn)
		}
	}
}

func collectMetric(sec int64, fn func() []*model.MetricValue) {
	hostname, err := g.Hostname()
	if err != nil {
		return
	}

	mvs := []*model.MetricValue{}
	ignoreMetrics := g.Config().IgnoreMetrics
	items := fn()
	if items == nil || len(items) == 0 {
		return
	}

	for _, mv := range items {
		if b, ok := ignoreMetrics[mv.Metric]; ok && b {
			continue
		} else {
			mvs = append(mvs, mv)
		}
	}

	now := time.Now().Unix()
	for j := 0; j < len(mvs); j++ {
		mvs[j].Step = sec
		if mvs[j].Endpoint == "" {
			mvs[j].Endpoint = hostname
		}
		mvs[j].Timestamp = now
	}

	if store.GetStoreStatus() {
		g.SendToTransfer(mvs)
		return
	}

	if len(mvs) > 0 {
		buf, err := json.Marshal(mvs)
		if err != nil || len(buf) == 0 {
			return
		}
		store := store.GetStore()
		store.Update(buf)
	}
}
