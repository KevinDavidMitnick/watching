package funcs

import (
	libvirt "github.com/libvirt/libvirt-go"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
)

func KvmHostInfo(dom *libvirt.Domain) map[string]interface{} {
	host := g.Config().Default.HostName
	M := make(map[string]interface{})
	info := make(map[string]interface{})
	name, err := dom.GetName()
	if err == nil && name != "" {
		info["hostname"] = name
	}

	platform, err := dom.GetOSType()
	if err == nil && platform != "" {
		info["platform"] = platform
	}

	proc, err := util.GetProcStat(dom)
	if err == nil && proc != nil && len(proc) == 1 {
		info["procs"] = proc[0]
	}

	status, err := util.GetOperStatus(dom)
	if err == nil && status != "" {
		info["status"] = status
	} else {
		info["status"] = "unknown"
	}

	info["NodeID"] = host
	M["hostInfo"] = info
	return M
}
