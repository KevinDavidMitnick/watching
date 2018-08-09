package g

import (
	"bytes"
	"encoding/json"
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

func getRrdLeader(addrs []string) string {
	var clusterStat RrdClusterStat
	for _, addr := range addrs {
		url := "http://" + addr + "/status"
		if resp, err := http.Get(url); err == nil {
			defer resp.Body.Close()
			if err1 := json.NewDecoder(resp.Body).Decode(&clusterStat); err1 == nil {
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
		log.Println("get rrd leader failed...")
		return nil
	}
	url = "http://" + url + "/db/execute?pretty&timings"
	if b, err := json.Marshal(data); err == nil && url != "" {
		log.Println("-----------------start flush------")
		//log.Println(string(b))
		resp, err := http.Post(url, "application/json", bytes.NewReader(b))
		if err != nil {
			log.Println("fail to flush", filename, len(items))
			return nil
		}
		defer resp.Body.Close()
		if _, err1 := ioutil.ReadAll(resp.Body); err1 == nil {
			log.Println("success to flush", filename, len(items))
			return err1
		}
	}
	return nil
}
