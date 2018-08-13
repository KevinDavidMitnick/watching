package g

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"io/ioutil"
	"net/http"
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

func postData(addr string, buf []byte) ([]byte, error) {
	log.Infof("send :%s,data :%s", addr, bytes.NewBuffer(buf).String())
	request, _ := http.NewRequest("POST", addr, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
	return nil, err
}

func getData(addr string) ([]byte, error) {
	log.Infof("send :%s", addr)
	request, _ := http.NewRequest("GET", addr, nil)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			return nil, fmt.Errorf("get response err.")
		}
		return ioutil.ReadAll(resp.Body)
	}
	return nil, err
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

	url := getRrdLeader(Config().Rrd.Addr)
	if url == "" {
		log.Errorln("get rrd leader failed...")
		return nil
	}
	url = "http://" + url + "/db/execute?pretty&timings"
	if b, err := json.Marshal(data); err == nil && url != "" {
		log.Infoln("-----------------start flush------")
		//log.Infoln(string(b))
		_, err := postData(url, b)
		if err != nil {
			log.Errorln("fail to flush", filename, len(items))
			return nil
		}
		log.Infoln("success to flush", filename, len(items))
	}
	return nil
}
