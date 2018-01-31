package util

import (
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"log"
	"time"
)

func GetCpuUtils(domain *libvirt.Domain) (map[string]float64, error) {
	cpuTime := make(map[string]uint64)
	cpuUtil := make(map[string]float64)
	vcpuInfo1, err := domain.GetVcpus()
	if err != nil {
		log.Println("get vcpus error", err)
		return cpuUtil, err
	}
	for _, cpuInfo := range vcpuInfo1 {
		core := fmt.Sprintf("cpucore=%d", cpuInfo.Number)
		cpuTime[core] = cpuInfo.CpuTime
	}
	startTime := time.Now()
	time.Sleep(time.Duration(1) * time.Second)
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	vcpuInfo2, err := domain.GetVcpus()
	for _, cpuInfo := range vcpuInfo2 {
		core := fmt.Sprintf("cpucore=%d", cpuInfo.Number)
		cpuUtil[core] = 100 * float64(cpuInfo.CpuTime-cpuTime[core]) / float64(duration)
	}
	return cpuUtil, nil
}
