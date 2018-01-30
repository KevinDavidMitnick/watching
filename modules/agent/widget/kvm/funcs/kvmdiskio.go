package funcs

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
	"time"
)

func KvmDiskIOMetrics() (L []*model.MetricValue) {
	host := g.Config().Default.HttpHost
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
		disks, err := util.GetDiskIO(&dom)
		if err != nil || disks == nil {
			continue
		}

		for k, v := range disks {
			metric := NewMetricValue(time.Now().Unix(), uuid, "disk.io.read_bytes", v[0], "GAUGE", "domain="+domain, k)
			L = append(L, metric)

			metric = NewMetricValue(time.Now().Unix(), uuid, "disk.io.read_requests", v[1], "GAUGE", "domain="+domain, k)
			L = append(L, metric)

			metric = NewMetricValue(time.Now().Unix(), uuid, "disk.io.write_bytes", v[2], "GAUGE", "domain="+domain, k)
			L = append(L, metric)

			metric = NewMetricValue(time.Now().Unix(), uuid, "disk.io.write_requests", v[3], "GAUGE", "domain="+domain, k)
			L = append(L, metric)

			metric = NewMetricValue(time.Now().Unix(), uuid, "disk.io.errs", v[4], "GAUGE", "domain="+domain, k)
			L = append(L, metric)
		}
		dom.Free()
	}
	return L
}
