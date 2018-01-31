package util

import (
	"encoding/json"
	libvirt "github.com/libvirt/libvirt-go"
)

type SwapStat struct {
	Available float64 `json:"available"`
	Out       float64 `json:"out"`
	Used      float64 `json:"used"`
	In        float64 `json:"in"`
	Total     float64 `json:"total"`
}

func GetSwapStat(domain *libvirt.Domain) ([]int64, error) {
	var swaps []int64
	uuid, err := domain.GetUUIDString()
	if err != nil {
		return nil, err
	}
	swapStat, err := GetValueFromQga(uuid, "guest-get-swap-info")
	if swapStat == nil || err != nil {
		return nil, err
	}
	var stat []SwapStat
	data, err := json.Marshal(swapStat)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &stat)
	for _, swap := range stat {
		total := swap.Total
		if int64(total) == 0 {
			swaps = []int64{0, 0}
			break
		}
		available := swap.Available
		used := total - available
		if int64(swap.Used) != 0 || int64(swap.In) != 0 || int64(swap.Out) != 0 {
			used = used * 1024
		}
		usage := (total - available) * 100 / total
		swaps = []int64{int64(used), int64(usage)}
	}
	return swaps, nil
}
