package cron

import (
	"encoding/json"
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	kvmfuncs "github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/funcs"
	kvmutil "github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
	"github.com/toolkits/net/httplib"
	"log"
	"strings"
	"time"
)

type KvmAndInterval struct {
	Fs       []func(*libvirt.Domain) map[string]interface{}
	Interval int
}

func ReportKvmInfo() {
	if g.Config().Kvm.Enabled && g.Config().Consul.Enabled && g.Config().Consul.Addr != "" && g.Config().Kvm.Host != "" {
		interval := g.Config().Transfer.Interval
		kvm := []KvmAndInterval{
			KvmAndInterval{
				Fs: []func(*libvirt.Domain) map[string]interface{}{
					kvmfuncs.KvmHostInfo,
					kvmfuncs.KvmCpuInfo,
					kvmutil.KvmDiskInfo,
					kvmutil.GetInterfaceInfo,
					kvmutil.KvmLoadInfo,
					kvmutil.KvmMemInfo,
				},
				Interval: interval,
			},
		}

		for _, v := range kvm {
			go reportKvmInfo(int64(v.Interval), v.Fs)
		}
	}
}

func reportKvmInfo(sec int64, fns []func(*libvirt.Domain) map[string]interface{}) {
	for {
		go MetricToConsul(sec, fns)
		time.Sleep(time.Duration(sec) * time.Second)
	}
}

func MetricToConsul(sec int64, fns []func(*libvirt.Domain) map[string]interface{}) {
	host := g.Config().Kvm.Host
	consulAddr := g.Config().Consul.Addr
	conn, err := kvmutil.ConnectLibvirt(host)
	if err != nil || conn == nil {
		return
	}

	doms, err := kvmutil.ListAliveDomains(conn)
	if err != nil {
		return
	}

	for _, dom := range doms {
		all := "{\"auto\":{%s}}"
		var gathers []string
		active, err := dom.IsActive()
		if active != true || err != nil {
			continue
		}

		uuid, err := dom.GetUUIDString()
		if err != nil || uuid == "" {
			continue
		}

		for _, fn := range fns {
			items := fn(&dom)
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}
			item, err := json.Marshal(items)
			if err == nil {
				gathers = append(gathers, strings.TrimPrefix(strings.TrimSuffix(string(item), "}"), "{"))
			}

		}
		vms := fmt.Sprintf(all, strings.Join(gathers, ","))
		reportKvmStatus(uuid, uuid, g.VERSION)
		sendKvmInfoToConsul(consulAddr, uuid, vms)
		dom.Free()
	}
	conn.Close()
}

// add by liucong.
func reportKvmStatus(hostname string, ip string, version string) {
	req := model.AgentReportRequest{
		Hostname:      hostname,
		IP:            ip,
		AgentVersion:  version,
		PluginVersion: "",
	}

	var resp model.SimpleRpcResponse
	err := g.HbsClient.Call("Agent.ReportStatus", req, &resp)
	if err != nil || resp.Code != 0 {
		log.Println("call Agent.ReportStatus fail:", err, "Request:", req, "Response:", resp)
	}
}

func sendKvmInfoToConsul(consulAddr string, uuid string, data string) {
	if data != "" {
		url := fmt.Sprintf("%s/v1/kv/kvm/%s?raw", consulAddr, uuid)
		r := httplib.Put(url)
		r.Body(data)
		ret, err := r.String()
		fmt.Println(url, data, ret, err)
	}
}
