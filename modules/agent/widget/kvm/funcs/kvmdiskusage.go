package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
	"time"
)

func KvmDiskUsageMetrics() (L []*model.MetricValue) {
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
		disks, err := util.GetDiskUsage(&dom)
		if err != nil || disks == nil {
			continue
		}

		for k, v := range disks {
			metric := NewMetricValue(time.Now().Unix(), uuid, "df.bytes.total", v[0], "GAUGE", "domain="+domain, k)
			L = append(L, metric)

			metric = NewMetricValue(time.Now().Unix(), uuid, "df.bytes.free", v[1], "GAUGE", "domain="+domain, k)
			L = append(L, metric)

			metric = NewMetricValue(time.Now().Unix(), uuid, "disk.bytes.used.percent", v[2], "GAUGE", "domain="+domain, k)
			L = append(L, metric)
		}
		dom.Free()
	}
	return L
}
