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

package http

import (
	"encoding/json"
	"fmt"
	"github.com/toolkits/net/httplib"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	cmodel "github.com/open-falcon/falcon-plus/common/model"

	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/gateway/g"
	trpc "github.com/open-falcon/falcon-plus/modules/gateway/receiver/rpc"
)

func configApiHttpRoutes() {
	http.HandleFunc("/api/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "blank body", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var metrics []*cmodel.MetricValue
		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "decode error", http.StatusBadRequest)
			return
		}

		reply := &cmodel.TransferResponse{}
		trpc.RecvMetricValues(metrics, reply, "http")

		RenderDataJson(w, reply)
	})
}

func configKvHttpRoutes() {
	http.HandleFunc("/kv/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "blank body", http.StatusBadRequest)
			return
		}

		if !g.Config().Consul.Enabled {
			http.Error(w, "not enable consul in gateway", http.StatusBadRequest)
			return
		}

		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "read data error", http.StatusBadRequest)
			return
		}
		err = req.ParseForm()
		if err != nil {
			http.Error(w, "param error", http.StatusBadRequest)
			return
		}
		url := req.Form.Get("url")
		if !strings.Contains(url, "/v1/kv") {
			http.Error(w, "param pattern error", http.StatusBadRequest)
			return
		}
		url = g.Config().Consul.Addr + url + "?raw"

		r := httplib.Put(url).SetTimeout(5*time.Second, 30*time.Second)
		r.Body(data)
		ret, _ := r.String()
		RenderDataJson(w, ret)
	})
}

func configGraphHttpRoutes() {
	http.HandleFunc("/graph/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "blank body", http.StatusBadRequest)
			return
		}

		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "read data error", http.StatusBadRequest)
			return
		}

		if len(data) == 0 {
			http.Error(w, "req data is empty", http.StatusBadRequest)
			return
		}

		ret := pushToRrd(data)
		RenderDataJson(w, ret)
	})
}

func pushToRrd(data []byte) string {
	items := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(data, &items); err != nil {
		return "unmarsh data failed."
	}

	filename := ""
	graphItems := make([]*cmodel.GraphItem, len(items))
	for i, item := range items {
		var graphItem cmodel.GraphItem
		graphItem.Endpoint = item["endpoint"].(string)
		graphItem.DsType = item["counterType"].(string)
		graphItem.Step = int(item["step"].(float64))
		graphItem.Metric = item["metric"].(string)
		graphItem.Timestamp = int64(item["timestamp"].(float64))
		graphItem.Value = item["value"].(float64)
		graphItem.Tags = utils.DictedTagstring(item["tags"].(string))
		if i == 0 {
			checksum := utils.Checksum(graphItem.Endpoint, graphItem.Metric, graphItem.Tags)

			filename = fmt.Sprintf("%s/%s_%s_%d.rrd", checksum[0:2], checksum, graphItem.DsType, graphItem.Step)
		}
		graphItems = append(graphItems, &graphItem)
	}
	if filename == "" {
		return "push to rrd data is empty,filename not right."
	}
	if err := g.Flushrrd(filename, graphItems); err == nil {
		return "success push data to rrd file:" + filename
	} else {
		return "failed to push data to rrd file:" + filename
	}
}
