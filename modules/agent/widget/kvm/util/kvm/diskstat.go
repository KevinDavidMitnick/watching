package util

import (
	"encoding/json"
	"fmt"
	"github.com/beevik/etree"
	libvirt "github.com/libvirt/libvirt-go"
	"strings"
	"time"
)

type DiskUsage struct {
	UsedKiloBytes  float64 `json:"used_kilobytes"`
	FreeKiloBytes  float64 `json:"free_kilobytes"`
	TotalKiloBytes float64 `json:"total_kilobytes"`
	Partition      string  `json:"partition"`
	MountedOn      string  `json:"mounted_on"`
}

func GetDiskIO(domain *libvirt.Domain) (map[string][]int64, error) {
	devs := make(map[string][]int64)
	xml, err := domain.GetXMLDesc(0)
	if err != nil {
		return devs, err
	}

	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return devs, err
	}

	for _, elem := range doc.FindElements("//devices/disk[@device='disk']/target") {
		for _, attr := range elem.Attr {
			if attr.Key == "dev" {
				blockStatsStart, err := domain.BlockStats(attr.Value)
				if err != nil {
					return devs, err
				}
				time.Sleep(time.Duration(1) * time.Second)
				blockStatsEnd, err := domain.BlockStats(attr.Value)
				if err != nil {
					return devs, err
				}
				rdBytes := int64(blockStatsEnd.RdBytes - blockStatsStart.RdBytes)
				rdReq := int64(blockStatsEnd.RdReq - blockStatsStart.RdReq)
				wrBytes := int64(blockStatsEnd.WrBytes - blockStatsStart.WrBytes)
				wrReq := int64(blockStatsEnd.WrReq - blockStatsStart.WrReq)
				errs := int64(blockStatsEnd.Errs - blockStatsStart.Errs)
				devs["device="+attr.Value] = []int64{rdBytes, rdReq, wrBytes, wrReq, errs}
			}
		}
	}
	return devs, nil
}

func GetDiskUsage(domain *libvirt.Domain) (map[string][]int64, error) {
	devs := make(map[string][]int64)
	uuid, err := domain.GetUUIDString()
	if err != nil {
		return nil, err
	}
	diskUsage, err := GetValueFromQga(uuid, "guest-get-disk-usage")
	if diskUsage == nil || err != nil {
		return nil, err
	}
	var usage []DiskUsage
	us, err := json.Marshal(diskUsage)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(us, &usage)
	for _, disk := range usage {
		key := fmt.Sprintf("fstype=%s,mount=%s", disk.Partition, disk.MountedOn)
		var usePercent float64 = 0
		if disk.TotalKiloBytes != 0 && disk.UsedKiloBytes != 0 {
			usePercent = 100 * disk.UsedKiloBytes / disk.TotalKiloBytes
		}
		devs[key] = []int64{int64(disk.TotalKiloBytes), int64(disk.FreeKiloBytes), int64(usePercent)}
	}
	return devs, nil
}

func KvmDiskInfo(dom *libvirt.Domain) map[string]interface{} {
	M := make(map[string]interface{})
	info := make(map[string]interface{})

	diskUsage, err := GetDiskUsage(dom)
	if err == nil && diskUsage != nil {
		for partition, usage := range diskUsage {
			arr := strings.Split(partition, ",")
			fstype := strings.Split(arr[0], "=")[1]
			mount := strings.Split(arr[1], "=")[1]
			if mount == "" {
				mount = fstype
			}
			info[mount] = map[string]interface{}{"fstype": fstype, "mountpoint": mount,
				"total": usage[0] * 1024, "free": usage[1] * 1024, "usedPercent": usage[2] * 1024,
			}
		}
	}

	M["diskInfo"] = info
	return M
}
