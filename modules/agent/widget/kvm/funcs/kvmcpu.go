package funcs

import (
	libvirt "github.com/libvirt/libvirt-go"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
	"time"
)

func KvmCpuMetrics() (L []*model.MetricValue) {
	host := g.Config().Kvm.Host
	conn, err := util.ConnectLibvirt(host)
	if err == nil && conn != nil {
		defer conn.Close()
	} else {
		return L
	}

	doms, err := util.ListAliveDomains(conn)
	if err != nil {
		return L
	}

	for _, dom := range doms {
		active, err := dom.IsActive()
		if active != true || err != nil {
			continue
		}

		uuid, err := dom.GetUUIDString()
		if err != nil {
			continue
		}
		domain, _ := dom.GetName()
		utils, err := util.GetCpuUtils(&dom)
		if err != nil || utils == nil {
			continue
		}

		for k, v := range utils {
			metric := NewMetricValue(time.Now().Unix(), uuid, "cpu.used.percent", v, "GAUGE", "domain="+domain, k)
			L = append(L, metric)
		}

		dom.Free()
	}
	return L
}

func KvmCpuInfo(dom *libvirt.Domain) map[string]interface{} {
	M := make(map[string]interface{})
	info := make(map[string]interface{})

	utils, err := util.GetCpuUtils(dom)
	coreNumber := 0
	usedPercent := 0.0

	if err == nil && utils != nil {
		for _, percent := range utils {
			coreNumber++
			usedPercent += percent
		}
		if coreNumber != 0 && usedPercent != 0.0 {
			info["coreNumber"] = coreNumber
			info["usedPercent"] = usedPercent / float64(coreNumber)
		}
	}

	M["cpuInfo"] = info
	return M
}
