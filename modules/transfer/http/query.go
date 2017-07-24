package http

import (
	"encoding/json"
	"fmt"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	"net/http"
)

type QueryMetric struct {
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
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

	decoder := json.NewDecoder(req.Body)
	var metrics []*QueryMetric
	err := decoder.Decode(&metrics)
	if err != nil || len(metrics) != 1 {
		http.Error(rw, "decode error", http.StatusBadRequest)
		return
	}

	metric := metrics[0].Metric
	endpoint := metrics[0].Endpoint
	tags := metrics[0].Tags
	str := ""
	for key, value := range tags {
		if str != "" {
			str = fmt.Sprintf("%s%s", str, ",")
		}
		str = fmt.Sprintf("%s%s=%s", str, key, value)
	}
	str = fmt.Sprintf("%s%s=%s", str, "endpoint", endpoint)
	query_str := fmt.Sprintf("timeseries=%s{%s}", metric, str)

	url := fmt.Sprintf("http://%s/api/query/last?%s", g.Config().Tsdb.Address, query_str)
	reply, err := http.Get(url)

	if err != nil {
		http.Error(rw, "tsdb query error", http.StatusBadRequest)
		return
	}

	RenderDataJson(rw, reply)
}

func configQueryRoutes() {
	http.HandleFunc("/query/last", query_last_metric)
}
