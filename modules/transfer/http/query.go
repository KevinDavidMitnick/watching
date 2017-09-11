// add by liucong,2017.7.25
package http

import (
	"encoding/json"
	"fmt"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"io/ioutil"
	"net/http"
)

type Metrics struct {
	Endpoint string            `json:"endpoint"`
	Metric   string            `json:"metric"`
	Tags     map[string]string `json:"tags"`
}

func query_last_metric(rw http.ResponseWriter, req *http.Request) {
	if !g.Config().Tsdb.Enabled {
		http.Error(rw, "tsdb is not enable", http.StatusBadRequest)
		return
	}

	if req.ContentLength == 0 {
		http.Error(rw, "blank body", http.StatusBadRequest)
		return
	}

	var metrics Metrics
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "read error", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &metrics)
	if err != nil {
		http.Error(rw, "decode error", http.StatusBadRequest)
		return
	}

	metric := metrics.Metric
	endpoint := metrics.Endpoint
	if err != nil {
		http.Error(rw, "decode error", http.StatusBadRequest)
		return
	}
	str := ""
	for key, value := range metrics.Tags {
		if str != "" {
			str = fmt.Sprintf("%s%s", str, ",")
		}
		str = fmt.Sprintf("%s%s=%s", str, key, value)
	}
	if str != "" {
		str = fmt.Sprintf("%s%s", str, ",")
	}
	str = fmt.Sprintf("%s%s=%s", str, "endpoint", endpoint)
	query_str := fmt.Sprintf("timeseries=%s{%s}", metric, str)

	url := fmt.Sprintf("http://%s/api/query/last?%s", g.Config().Tsdb.Address, query_str)
	reply, err := http.Get(url)

	if err != nil {
		http.Error(rw, "tsdb query error", http.StatusBadRequest)
		return
	}

	body, err = ioutil.ReadAll(reply.Body)
	if err != nil {
		http.Error(rw, "request to tsdb failed.", http.StatusBadRequest)
		return
	}

	var ret []map[string]interface{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		http.Error(rw, "decode error", http.StatusBadRequest)
		return
	}

	RenderDataJson(rw, ret)
}

func configQueryRoutes() {
	http.HandleFunc("/query/last", query_last_metric)
}
