package util

import (
	"encoding/json"
	libvirt "github.com/libvirt/libvirt-go"
)

func GetProcStat(domain *libvirt.Domain) ([]int64, error) {
	var procs []int64
	uuid, err := domain.GetUUIDString()
	if err != nil {
		return nil, err
	}
	procStat, err := GetValueFromQga(uuid, "guest-get-process-number")
	if procStat == nil || err != nil {
		return nil, err
	}
	var stat float64
	data, err := json.Marshal(procStat)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &stat)
	procs = []int64{int64(stat)}
	return procs, nil
}
