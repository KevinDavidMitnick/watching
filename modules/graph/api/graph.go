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

package api

import (
	//"fmt"
	//"math"
	log "github.com/Sirupsen/logrus"
	pfc "github.com/niean/goperfcounter"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
	"github.com/open-falcon/falcon-plus/modules/graph/rrdtool"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

type Graph int

func (this *Graph) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (this *Graph) Send(items []*cmodel.GraphItem, resp *cmodel.SimpleRpcResponse) error {
	go handleItems(items)
	return nil
}

// 供外部调用、处理接收到的数据 的接口
func HandleItems(items []*cmodel.GraphItem) error {
	handleItems(items)
	return nil
}

// func handleItems(items []*cmodel.GraphItem) {
// 	if items == nil {
// 		return
// 	}

// 	count := len(items)
// 	if count == 0 {
// 		return
// 	}

// 	cfg := g.Config()

// 	for i := 0; i < count; i++ {
// 		if items[i] == nil {
// 			continue
// 		}

// 		endpoint := items[i].Endpoint
// 		if !g.IsValidString(endpoint) {
// 			if cfg.Debug {
// 				log.Printf("invalid endpoint: %s", endpoint)
// 			}
// 			pfc.Meter("invalidEnpoint", 1)
// 			continue
// 		}

// 		counter := cutils.Counter(items[i].Metric, items[i].Tags)
// 		if !g.IsValidString(counter) {
// 			if cfg.Debug {
// 				log.Printf("invalid counter: %s/%s", endpoint, counter)
// 			}
// 			pfc.Meter("invalidCounter", 1)
// 			continue
// 		}

// 		dsType := items[i].DsType
// 		step := items[i].Step
// 		checksum := items[i].Checksum()
// 		key := g.FormRrdCacheKey(checksum, dsType, step)

// 		//statistics
// 		proc.GraphRpcRecvCnt.Incr()

// 		// To Graph
// 		first := store.GraphItems.First(key)
// 		if first != nil && items[i].Timestamp <= first.Timestamp {
// 			continue
// 		}
// 		store.GraphItems.PushFront(key, items[i], checksum, cfg)

// 		// To Index
// 		index.ReceiveItem(items[i], checksum)

// 		// To History
// 		store.AddItem(checksum, items[i])
// 	}
// }

func handleItems(items []*cmodel.GraphItem) {
	if items == nil {
		return
	}

	count := len(items)
	if count == 0 {
		return
	}

	cfg := g.Config()

	for i := 0; i < count; i++ {
		if items[i] == nil {
			continue
		}

		endpoint := items[i].Endpoint
		if !g.IsValidString(endpoint) {
			if cfg.Debug {
				log.Printf("invalid endpoint: %s", endpoint)
			}
			pfc.Meter("invalidEnpoint", 1)
			continue
		}

		counter := cutils.Counter(items[i].Metric, items[i].Tags)
		if !g.IsValidString(counter) {
			if cfg.Debug {
				log.Printf("invalid counter: %s/%s", endpoint, counter)
			}
			pfc.Meter("invalidCounter", 1)
			continue
		}

		dsType := items[i].DsType
		step := items[i].Step
		checksum := items[i].Checksum()
		filename := g.RrdFileName(checksum, dsType, step)

		//statistics
		proc.GraphRpcRecvCnt.Incr()

		// To Graph
		//first := store.GraphItems.First(key)
		//if first != nil && items[i].Timestamp <= first.Timestamp {
		//	continue
		//}
		//store.GraphItems.PushFront(key, items[i], checksum, cfg)

		// To Index
		index.ReceiveItem(items[i], checksum)

		// To History
		store.AddItem(checksum, items[i])

		// flush to disk
		rrdtool.Flushrrd(filename, []*cmodel.GraphItem{items[i]})
	}
}

func (this *Graph) Query(param cmodel.GraphQueryParam, resp *cmodel.GraphQueryResponse) error {
	var (
		datas       []*cmodel.RRDData
	)

	// statistics
	proc.GraphQueryCnt.Incr()

	// form empty response
	resp.Values = []*cmodel.RRDData{}
	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	dsType, step, exists := index.GetTypeAndStep(param.Endpoint, param.Counter) // complete dsType and step
	if !exists {
		return nil
	}
	resp.DsType = dsType
	resp.Step = step

	start_ts := param.Start - param.Start%int64(step)
	end_ts := param.End - param.End%int64(step) + int64(step)
	if end_ts-start_ts-int64(step) < 1 {
		return nil
	}

	md5 := cutils.Md5(param.Endpoint + "/" + param.Counter)
	filename := g.RrdFileName(md5, dsType, step)

	// read data from rrd file
	// 从RRD中获取数据不包含起始时间点
	// 例: start_ts=1484651400,step=60,则第一个数据时间为1484651460)
	datas, _ = rrdtool.Fetch(filename, param.ConsolFun, start_ts-int64(step), end_ts, step)
	//fmt.Println(datas)

	resp.Values = datas
	// statistics
	proc.GraphQueryItemCnt.IncrBy(int64(len(resp.Values)))
	return nil
}

//从内存索引、MySQL中删除counter，并从磁盘上删除对应rrd文件
func (this *Graph) Delete(params []*cmodel.GraphDeleteParam, resp *cmodel.GraphDeleteResp) error {
	resp = &cmodel.GraphDeleteResp{}
	for _, param := range params {
		err, tags := cutils.SplitTagsString(param.Tags)
		if err != nil {
			log.Error("invalid tags:", param.Tags, "error:", err)
			continue
		}

		var item *cmodel.GraphItem = &cmodel.GraphItem{
			Endpoint: param.Endpoint,
			Metric:   param.Metric,
			Tags:     tags,
			DsType:   param.DsType,
			Step:     param.Step,
		}
		index.RemoveItem(item)
	}

	return nil
}

func (this *Graph) Info(param cmodel.GraphInfoParam, resp *cmodel.GraphInfoResp) error {
	// statistics
	proc.GraphInfoCnt.Incr()

	dsType, step, exists := index.GetTypeAndStep(param.Endpoint, param.Counter)
	if !exists {
		return nil
	}

	//md5 := cutils.Md5(param.Endpoint + "/" + param.Counter)
	//filename := fmt.Sprintf("%s/%s/%s_%s_%d.rrd", g.Config().RRD.Storage, md5[0:2], md5, dsType, step)
	filename := ""

	resp.ConsolFun = dsType
	resp.Step = step
	resp.Filename = filename

	return nil
}

func (this *Graph) Last(param cmodel.GraphLastParam, resp *cmodel.GraphLastResp) error {
	// statistics
	proc.GraphLastCnt.Incr()

	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	resp.Value = GetLast(param.Endpoint, param.Counter)

	return nil
}

func (this *Graph) LastRaw(param cmodel.GraphLastParam, resp *cmodel.GraphLastResp) error {
	// statistics
	proc.GraphLastRawCnt.Incr()

	resp.Endpoint = param.Endpoint
	resp.Counter = param.Counter
	resp.Value = GetLastRaw(param.Endpoint, param.Counter)

	return nil
}

// 非法值: ts=0,value无意义
func GetLast(endpoint, counter string) *cmodel.RRDData {
	dsType, step, exists := index.GetTypeAndStep(endpoint, counter)
	if !exists {
		return cmodel.NewRRDData(0, 0.0)
	}

	if dsType == g.GAUGE {
		return GetLastRaw(endpoint, counter)
	}

	if dsType == g.COUNTER || dsType == g.DERIVE {
		md5 := cutils.Md5(endpoint + "/" + counter)
		items := store.GetAllItems(md5)
		if len(items) < 2 {
			return cmodel.NewRRDData(0, 0.0)
		}

		f0 := items[0]
		f1 := items[1]
		delta_ts := f0.Timestamp - f1.Timestamp
		delta_v := f0.Value - f1.Value
		if delta_ts != int64(step) || delta_ts <= 0 {
			return cmodel.NewRRDData(0, 0.0)
		}
		if delta_v < 0 {
			// when cnt restarted, new cnt value would be zero, so fix it here
			delta_v = 0
		}

		return cmodel.NewRRDData(f0.Timestamp, delta_v/float64(delta_ts))
	}

	return cmodel.NewRRDData(0, 0.0)
}

// 非法值: ts=0,value无意义
func GetLastRaw(endpoint, counter string) *cmodel.RRDData {
	md5 := cutils.Md5(endpoint + "/" + counter)
	item := store.GetLastItem(md5)
	return cmodel.NewRRDData(item.Timestamp, item.Value)
}
