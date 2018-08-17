package g

import (
	"bytes"
	"encoding/json"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/net/httplib"
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
	Items  []*cmodel.GraphItem `json:"items"`
	Method string              `json:"method"`
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

func getRrdLeader_old(addrs []string) string {
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

func getRrdLeader(addrs []string) string {
	var clusterStat RrdClusterStat
	for _, addr := range addrs {
		url := "http://" + addr + "/status"
		req := httplib.Get(url).SetTimeout(2*time.Second, 10*time.Second)
		err := req.ToJson(&clusterStat)
		if err == nil {
			if clusterStat.Store.Raft.State == "Leader" {
				return addr
			}
		}
		log.Errorln("### Get Leader data err: ", err)
	}
	return ""
}

// flush to disk from memory
// 最新的数据在列表的最后面
// TODO fix me, filename fmt from item[0], it's hard to keep consistent
func Flushrrd_old(filename string, items []*cmodel.GraphItem) error {
	var data Flushfile_t
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

func Flushrrd(items []*cmodel.GraphItem) error {
	var data Flushfile_t
	data.Items = items
	data.Method = "insert"

	url := getRrdLeader(Config().Rrd.Addr)
	if url == "" {
		log.Errorln("get rrd leader failed...")
		return nil
	}
	url = "http://" + url + "/db/execute?pretty&timings"
	req := httplib.Post(url).SetTimeout(3*time.Second, 10*time.Second)
	if b, err := json.Marshal(data); err == nil && url != "" {
		log.Infoln("-----------------start flush------")
		req.Body(b)
		_, err := req.String()
		if err != nil {
			log.Errorln("fail to flush", len(items))
			return nil
		}
		log.Infoln("success to flush", len(items))
	}
	return nil
}
