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
	"fmt"
	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
	"io/ioutil"
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
	Method   string `json:"method"`
}

type flushfile_t struct {
	filename string
	items    []*cmodel.GraphItem
}

type Flushfile_t struct {
	Filename string              `json:"filename"`
	Items    []*cmodel.GraphItem `json:"items"`
	Method   string              `json:"method"`
}

type Fetch_return struct {
	Result []*cmodel.RRDData `json:"results"`
	Time   float64           `json:"time"`
}

type RaftStat struct {
	State string `json:"state"`
}

type StoreStat struct {
	Raft *RaftStat `json:"raft"`
}

type RrdClusterStat struct {
	Store *StoreStat `json:"store"`
}

func Start() {
	// sync disk
	go syncDisk()
	go ioWorker()
	log.Println("rrdtool.Start ok")
}

func GetRrdLeader() string {
	addrs := g.Config().Rrd.Addr
	return getRrdLeader(addrs)
}

func getData(url string) ([]byte, error) {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("TIMEOUT", "10")
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Close = true

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func submitData(url string, data []byte, method string) ([]byte, error) {
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("TIMEOUT", "10")
	request.Close = true

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getRrdLeader(addrs []string) string {
	var clusterStat RrdClusterStat
	for _, addr := range addrs {
		url := "http://" + addr + "/status"
		if resp, err := getData(url); err == nil {
			if err1 := json.Unmarshal(resp, &clusterStat); err1 == nil {
				if clusterStat.Store.Raft.State == "Leader" {
					return addr
				}
			}
		}
	}
	return ""
}

// flush to disk from memory
// 最新的数据在列表的最后面
// TODO fix me, filename fmt from item[0], it's hard to keep consistent
func Flushrrd(filename string, items []*cmodel.GraphItem) error {
	var data Flushfile_t
	data.Filename = filename
	data.Items = items
	data.Method = "insert"

	url := getRrdLeader(g.Config().Rrd.Addr)
	if url == "" {
		log.Errorln("get rrd leader failed...")
		return nil
	}
	url = "http://" + url + "/db/execute?pretty&timings"
	if b, err := json.Marshal(data); err == nil && url != "" {
		log.Infoln("-----------------start flush------")
		//log.Infoln(string(b))
		_, err := submitData(url, b,"POST")
		if err != nil {
			log.Errorln("fail to flush", filename, len(items))
			return nil
		}
		log.Infoln("success to flush", filename, len(items))
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

// func Fetch(filename string, cf string, start, end int64, step int) ([]*cmodel.RRDData, error) {
// 	done := make(chan error, 1)
// 	task := &io_task_t{
// 		method: IO_TASK_M_FETCH,
// 		args: &fetch_t{
// 			filename: filename,
// 			cf:       cf,
// 			start:    start,
// 			end:      end,
// 			step:     step,
// 		},
// 		done: done,
// 	}
// 	io_task_chan <- task
// 	err := <-done
// 	return task.args.(*fetch_t).data, err
// }

func Fetch(filename string, cf string, start, end int64, step int) ([]*cmodel.RRDData, error) {
	var rrd []*cmodel.RRDData
	var fetch_return Fetch_return
	var data Fetch_t
	data.Start = start
	data.End = end
	data.Step = int(step)
	data.Cf = cf
	data.Filename = filename
	data.Method = "query"
	log.Println("starting fetching data....")
	if b, err := json.Marshal(data); err == nil {
		log.Println(string(b))
		url := getRrdLeader(g.Config().Rrd.Addr)
		if url == "" {
			log.Infoln("get rrd leader failed...")
			return rrd, nil
		}
		url = "http://" + url + "/db/query?pretty&timings"
		resp, err1 := submitData(url, b,"POST")
		if err1 != nil {
			log.Infof("fetch error:filename is %s,start time is:%d,end time is:%d,step is :%d,time_len is:%d", filename, start, end, step, len(rrd))
			return rrd, nil
		}
		if err3 := json.Unmarshal(resp, &fetch_return); err3 == nil {
			rrd = fetch_return.Result
			log.Infof("success fetch data,len is: %d", len(rrd))
			return rrd, nil
		}
	}
	return rrd, nil
}

func FlushAll(force bool) {
	n := store.GraphItems.Size / 10
	for i := 0; i < store.GraphItems.Size; i++ {
		FlushRRD(i, force)
		if i%n == 0 {
			log.Infof("flush hash idx:%03d size:%03d disk:%08d net:%08d\n",
				i, store.GraphItems.Size, disk_counter, net_counter)
		}
	}
	log.Infof("flush hash done (disk:%08d net:%08d)\n", disk_counter, net_counter)
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
