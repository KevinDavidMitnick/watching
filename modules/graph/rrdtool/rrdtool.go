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

package rrdtool

import (
	"bytes"
	"encoding/json"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	disk_counter     uint64
	net_counter      uint64
	flushrrd_timeout int32
)

type fetch_t struct {
	filename string
	cf       string
	start    int64
	end      int64
	step     int
	data     []*cmodel.RRDData
}

type Fetch_t struct {
	Filename string `json:"filename"`
	Cf       string `json:"cf"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Step     int    `json:"step"`
}

type flushfile_t struct {
	filename string
	items    []*cmodel.GraphItem
}

type Flushfile_t struct {
	Filename string              `json:"filename`
	Items    []*cmodel.GraphItem `json:items`
}

func Start() {
	// sync disk
	go syncDisk()
	go ioWorker()
	log.Println("rrdtool.Start ok")
}

// flush to disk from memory
// 最新的数据在列表的最后面
// TODO fix me, filename fmt from item[0], it's hard to keep consistent
func flushrrd(filename string, items []*cmodel.GraphItem) error {
	var data flushfile_t
	data.filename = filename
	data.items = items

	if b, err := json.Marshal(data); err == nil {
		url := g.Config().Rrd.AppendAddr
		resp, err := http.Post(url, "application/json", bytes.NewReader(b))
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		if ret, err1 := ioutil.ReadAll(resp.Body); err == nil {
			log.Println(filename, len(items), string(ret))
			return err1
		}
	}
	return nil
}

func FlushFile(filename string, items []*cmodel.GraphItem) error {
	done := make(chan error, 1)
	io_task_chan <- &io_task_t{
		method: IO_TASK_M_FLUSH,
		args: &flushfile_t{
			filename: filename,
			items:    items,
		},
		done: done,
	}
	atomic.AddUint64(&disk_counter, 1)
	return <-done
}

func Fetch(filename string, cf string, start, end int64, step int) ([]*cmodel.RRDData, error) {
	done := make(chan error, 1)
	task := &io_task_t{
		method: IO_TASK_M_FETCH,
		args: &fetch_t{
			filename: filename,
			cf:       cf,
			start:    start,
			end:      end,
			step:     step,
		},
		done: done,
	}
	io_task_chan <- task
	err := <-done
	return task.args.(*fetch_t).data, err
}

func fetch(filename string, cf string, start, end int64, step int) ([]*cmodel.RRDData, error) {
	var rrd []*cmodel.RRDData
	var data Fetch_t
	data.Start = start
	data.End = end
	data.Step = int(time.Duration(step) * time.Second)
	data.Cf = cf

	if b, err := json.Marshal(data); err == nil {
		log.Println(string(b))
		url := g.Config().Rrd.QueryAddr
		resp, err1 := http.Post(url, "application/json", bytes.NewReader(b))
		if err1 != nil {
			for i := 1; i <= 60; i++ {
				ts := i * step
				rrd = append(rrd, &cmodel.RRDData{Timestamp: start + int64(ts), Value: 3.0})
			}
			log.Println("start time is:%d,end time is:%d,step is :%d,time_len is:%d", start, end, step, len(rrd))
			return rrd, nil
		}
		defer resp.Body.Close()
		if ret, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			json.Unmarshal(ret, &rrd)
			log.Printf("fetch dat: %v", rrd)
		}
	}
	return rrd, nil
}

func FlushAll(force bool) {
	n := store.GraphItems.Size / 10
	for i := 0; i < store.GraphItems.Size; i++ {
		FlushRRD(i, force)
		if i%n == 0 {
			log.Printf("flush hash idx:%03d size:%03d disk:%08d net:%08d\n",
				i, store.GraphItems.Size, disk_counter, net_counter)
		}
	}
	log.Printf("flush hash done (disk:%08d net:%08d)\n", disk_counter, net_counter)
}

func CommitByKey(key string) {
	md5, dsType, step, err := g.SplitRrdCacheKey(key)
	if err != nil {
		return
	}
	filename := g.RrdFileName(md5, dsType, step)
	items := store.GraphItems.PopAll(key)
	if len(items) == 0 {
		return
	}
	FlushFile(filename, items)
}

func FlushRRD(idx int, force bool) {
	atomic.StoreInt32(&flushrrd_timeout, 0)
	keys := store.GraphItems.KeysByIndex(idx)
	if len(keys) == 0 {
		return
	}
	for _, key := range keys {
		if force || shouldFlush(key) {
			CommitByKey(key)
		}
	}
}

func shouldFlush(key string) bool {
	if store.GraphItems.ItemCnt(key) >= g.FLUSH_MIN_COUNT {
		return true
	}

	deadline := time.Now().Unix() - int64(g.FLUSH_MAX_WAIT)
	back := store.GraphItems.Back(key)
	if back != nil && back.Timestamp <= deadline {
		return true
	}

	return false
}
