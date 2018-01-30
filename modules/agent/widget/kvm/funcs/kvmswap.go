package funcs

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
	"time"
)

func KvmSwapMetrics() (L []*model.MetricValue) {
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
		swap, err := util.GetSwapStat(&dom)
		if err != nil || swap == nil || len(swap) != 2 {
			continue
		}

		metric := NewMetricValue(time.Now().Unix(), uuid, "mem.swapused", swap[0], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "mem.swapused.percent", swap[1], "GAUGE", "domain="+domain)
		L = append(L, metric)
		dom.Free()
	}
	return L
}
