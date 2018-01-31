package funcs

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
	"time"
)

func KvmTcpMetrics() (L []*model.MetricValue) {
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
		tcp, err := util.GetTcpStat(&dom)
		if err != nil || tcp == nil || len(tcp) != 5 {
			continue
		}

		metric := NewMetricValue(time.Now().Unix(), uuid, "ss.estab", tcp[0], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "ss.synrecv", tcp[1], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "ss.timewait", tcp[2], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "ss.closewait", tcp[3], "GAUGE", "domain="+domain)
		L = append(L, metric)

		metric = NewMetricValue(time.Now().Unix(), uuid, "ss.listen", tcp[4], "GAUGE", "domain="+domain)
		L = append(L, metric)
		dom.Free()
	}
	return L
}
