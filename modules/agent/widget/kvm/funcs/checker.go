package funcs

import (
	"fmt"
)

func CheckKvmCollector() {
	output := make(map[string]bool)

	output["CpuMetrics"] = len(KvmCpuMetrics()) > 0
	output["MemMetrics"] = len(KvmMemMetrics()) > 0
	output["NetMetrics"] = len(KvmNetMetrics()) > 0
	output["DiskIOMetrics"] = len(KvmDiskIOMetrics()) > 0
	output["DiskUsageMetrics"] = len(KvmDiskUsageMetrics()) > 0
	output["LoadMetrics"] = len(KvmLoadMetrics()) > 0
	output["ProcMetrics"] = len(KvmProcMetrics()) > 0
	output["SwapMetrics"] = len(KvmSwapMetrics()) > 0
	output["TcpMetrics"] = len(KvmTcpMetrics()) > 0

	for k, v := range output {
		status := "fail"
		if v {
			status = "ok"
		}
		fmt.Println(k, "...", status)
	}

}
