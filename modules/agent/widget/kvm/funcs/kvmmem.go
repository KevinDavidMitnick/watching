package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
	"time"
)

func KvmMemMetrics() (L []*model.MetricValue) {
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
		mem, err := util.GetMemStat(&dom)
		if err != nil || mem == nil || len(mem) != 4 {
			continue
		}

		metric := NewMetricValue(time.Now().Unix(), uuid, "mem.memtotal", mem[0], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "mem.memused", mem[1], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "mem.memfree", mem[2], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "mem.memused.percent", mem[3], "GAUGE", "domain="+domain)
		L = append(L, metric)

		dom.Free()
	}
	return L
}
