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
	"github.com/toolkits/net/httplib"
	"io/ioutil"
	"net/http"
	"strings"

	cmodel "github.com/open-falcon/falcon-plus/common/model"

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

		r := httplib.Put(url)
		r.Body(data)
		ret, _ := r.String()
		RenderDataJson(w, ret)
	})
}
