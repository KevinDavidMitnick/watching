package funcs

import (
	libvirt "github.com/libvirt/libvirt-go"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/open-falcon/falcon-plus/modules/agent/widget/kvm/util"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

type AbilityAndInterval struct {
	Fs       []func(*libvirt.Domain) map[string]interface{}
	Interval int
}

var Mappers []FuncsAndInterval
var Abilities []AbilityAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmCpuMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmMemMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmNetMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmDiskIOMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmDiskUsageMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmLoadMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmProcMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmSwapMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				KvmTcpMetrics,
			},
			Interval: interval,
		},
	}
}

func BuildAbilities() {
	interval := g.Config().Kvm.Interval
	Abilities = []AbilityAndInterval{
		AbilityAndInterval{
			Fs: []func(*libvirt.Domain) map[string]interface{}{
				KvmHostInfo,
				KvmCpuInfo,
				util.KvmDiskInfo,
				util.GetInterfaceInfo,
				util.KvmLoadInfo,
				util.KvmMemInfo,
			},
			Interval: interval,
		},
	}
}
