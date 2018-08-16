package g

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"io/ioutil"
	"net/http"
	"time"
)

type RaftStat struct {
	State string `json:"state"`
}

type StoreStat struct {
	Raft *RaftStat `json:"raft"`
}

type RrdClusterStat struct {
	Store *StoreStat `json:"store"`
}

type Flushfile_t struct {
	Filename string              `json:"filename"`
	Items    []*cmodel.GraphItem `json:"items"`
	Method   string              `json:"method"`
}

// GetData from url,use method get
func GetData(url string) ([]byte, error) {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Close = true

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(10) * time.Second,
	}
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

// SubmitData from url,use method submit
func SubmitData(url string, data []byte, method string) ([]byte, error) {
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Close = true

	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(10) * time.Second,
	}
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
		if resp, err := GetData(url); err == nil {
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

	url := getRrdLeader(Config().Rrd.Addr)
	if url == "" {
		log.Errorln("get rrd leader failed...")
		return nil
	}
	url = "http://" + url + "/db/execute?pretty&timings"
	if b, err := json.Marshal(data); err == nil && url != "" {
		log.Infoln("-----------------start flush------")
		//log.Infoln(string(b))
		_, err := SubmitData(url, b, "POST")
		if err != nil {
			log.Errorln("fail to flush", filename, len(items))
			return nil
		}
		log.Infoln("success to flush", filename, len(items), string(b))
	}
	return nil
}
