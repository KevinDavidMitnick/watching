package util

import (
	"fmt"
	"testing"
)

func TestDomain(t *testing.T) {
	conn, err := ConnectLibvirt("10.0.20.11")
	if err == nil && conn != nil {
		defer conn.Close()
	}
	//doms, _ := ListAliveDomains(conn)
	doms, _ := ListAliveDomains(conn)

	for _, dom := range doms {
		defer dom.Free()
		name, _ := dom.GetName()
		fmt.Println(name)
		// print cpu msg.
		//vcpuInfo, _ := dom.GetVcpus()
		//fmt.Println(len(vcpuInfo))
		//utils, _ := GetCpuUtils(&dom)

		//devs, _ := GetDiskIO(&dom)

		//devs, err := GetDiskUsage(&dom)

		//mems, err := GetMemStat(&dom)

		//swaps, err := GetSwapStat(&dom)

		//procs, err := GetProcStat(&dom)

		//tcps, err := GetTcpStat(&dom)

		// oper, err := GetOperStatus(&dom)

		// load, err := GetLoadStat(&dom)

		//inter, err := GetInterfaceIO(&dom)

		info, err := GetInterfaceInfo(&dom)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(info, err)
		}
		break
	}
}
