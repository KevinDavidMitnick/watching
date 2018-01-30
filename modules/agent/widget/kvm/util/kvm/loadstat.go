package util

import (
	"encoding/json"
	libvirt "github.com/libvirt/libvirt-go"
)

type LoadStat struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

func GetLoadStat(domain *libvirt.Domain) ([]int64, error) {
	var loads []int64
	uuid, err := domain.GetUUIDString()
	if err != nil {
		return nil, err
	}
	loadStat, err := GetValueFromQga(uuid, "guest-get-systemload-info")
	if loadStat == nil || err != nil {
		return nil, err
	}
	var stat []LoadStat
	data, err := json.Marshal(loadStat)
	if err != nil || data == nil || len(data) == 0 {
		return nil, err
	}
	json.Unmarshal(data, &stat)
	for _, load := range stat {
		load1 := load.Load1
		load5 := load.Load5
		load15 := load.Load15
		loads = []int64{int64(load1), int64(load5), int64(load15)}
	}
	return loads, nil
}

func KvmLoadInfo(dom *libvirt.Domain) map[string]interface{} {
	M := make(map[string]interface{})
	info := make(map[string]interface{})

	loadStat, err := GetLoadStat(dom)
	if err == nil && loadStat != nil && len(loadStat) == 3 {
		info["load1"] = loadStat[0]
		info["load5"] = loadStat[1]
		info["load15"] = loadStat[2]
	}
	M["loadInfo"] = info
	return M
}
