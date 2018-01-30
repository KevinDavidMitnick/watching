package util

import (
	"encoding/json"
	libvirt "github.com/libvirt/libvirt-go"
	"log"
)

type MemStat struct {
	Available float64 `json:"available"`
	Used      float64 `json:"used"`
	Cached    float64 `json:"cached"`
	Free      float64 `json:"free"`
	Total     float64 `json:"total"`
	Buffers   float64 `json:"buffers"`
}

func getMemUsage(domain *libvirt.Domain) ([]int64, error) {
	var total, used, available, usage int64 = 0, 0, 0, 0
	domainInfo, err := domain.GetInfo()
	if err != nil || domainInfo == nil {
		log.Println("get domain info error", err)
		return nil, err
	}
	total = int64(domainInfo.Memory) * 1024
	mem, _ := domain.MemoryStats(8, 0)
	if mem != nil && len(mem) > 1 {
		for _, memVal := range mem {
			switch libvirt.DomainMemoryStatTags(memVal.Tag) {
			case libvirt.DOMAIN_MEMORY_STAT_UNUSED:
				available = int64(memVal.Val * 1024)
				used = total - available
				usage = used * 100 / total
			}
		}
	}
	return []int64{int64(total), int64(used), int64(available), int64(usage)}, nil
}

func GetMemStat(domain *libvirt.Domain) ([]int64, error) {
	var mems []int64
	uuid, err := domain.GetUUIDString()
	if err != nil {
		return nil, err
	}
	memStat, err := GetValueFromQga(uuid, "guest-get-memory-info")
	if memStat == nil || err != nil {
		return getMemUsage(domain)
	}
	var stat []MemStat
	data, err := json.Marshal(memStat)
	if err != nil {
		return getMemUsage(domain)
	}
	err = json.Unmarshal(data, &stat)
	if err != nil {
		return getMemUsage(domain)
	}
	for _, mem := range stat {
		total := mem.Total
		available := mem.Available
		used := total - available
		usage := used * 100 / total
		mems = []int64{int64(total), int64(used), int64(available), int64(usage)}
	}
	return mems, nil
}

func KvmMemInfo(dom *libvirt.Domain) map[string]interface{} {
	M := make(map[string]interface{})
	info := make(map[string]interface{})

	memStat, err := GetMemStat(dom)
	if err == nil && memStat != nil && len(memStat) == 4 {
		info["total"] = memStat[0]
		info["used"] = memStat[1]
		info["available"] = memStat[2]
		info["usedPercent"] = memStat[3]
	}
	M["memInfo"] = info
	return M
}
