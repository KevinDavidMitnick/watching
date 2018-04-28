// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	gonet "net"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"errors"
	"io/ioutil"

	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/toolkits/net/httplib"
)

var (
	info map[string]interface{}
)

func ReportSystemInfo() {
	if g.Config().Consul.Enabled && g.Config().Consul.Addr != "" {
		go reportSystemInfo()
	}
}

func reportSystemInfo() {
	duration := time.Duration(g.Config().Consul.Interval) * time.Second
	consulUrl := g.Config().Consul.Addr
	nodeId := g.Config().Hostname

	for {
		time.Sleep(duration)
		data := GetAllInfo()
		sendSystemInfoToConsul(consulUrl, nodeId, data)
	}
}

func getHostInfo() map[string]interface{} {
	hostInfo, err := host.Info()
	if err == nil {
		uptime, _ := formatDeltatime(hostInfo.Uptime)
		host := map[string]interface{}{
			"hostname":             hostInfo.Hostname,
			"uptime":               uptime,
			"bootTime":             hostInfo.BootTime,
			"procs":                hostInfo.Procs,
			"os":                   hostInfo.OS,
			"platform":             hostInfo.Platform,
			"platformFamily":       hostInfo.PlatformFamily,
			"platformVersion":      hostInfo.PlatformVersion,
			"kernelVersion":        hostInfo.KernelVersion,
			"virtualizationSystem": hostInfo.VirtualizationSystem,
			"virtualizationRole":   hostInfo.VirtualizationRole,
			"hostid":               hostInfo.HostID,
		}
		info["hostInfo"] = host
		info["name"] = hostInfo.Hostname

		// add all disk info
		if os := hostInfo.OS; os != "" {
			switch os {
			case "linux":
				cmdStr := "df -x tmpfs -x devtmpfs --total|awk '/total/{print $2}'"
				cmd := exec.Command("/bin/sh", "-c", cmdStr)

				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					return info
				}

				diskStr := strings.TrimSpace(out.String())
				if diskTotal, err := strconv.ParseInt(diskStr, 10, 64); err == nil {
					info["diskTotal"] = diskTotal * 1024
				}
			}

		}

	}
	return info
}

func getLoadInfo() map[string]interface{} {
	loadInfo, err := load.Avg()
	if err == nil {
		info["loadInfo"] = loadInfo
	}
	return info
}

func getCpuInfo() map[string]interface{} {
	usedPercent, err1 := cpu.Percent(time.Second, false)
	coreNumber, err2 := cpu.Counts(false)
	cpuInfo, err3 := cpu.Info()
	if err1 == nil && err2 == nil && err3 == nil {
		info["cpuInfo"] = map[string]interface{}{
			"usedPercent": usedPercent[0],
			"coreNumber":  coreNumber,
			"modelName":   cpuInfo[0].ModelName,
		}
	}
	return info
}

func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		panic("ipv4Mask: len must be 4 bytes")
	}

	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

func getInterfaceInfo() map[string]interface{} {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	interfaceInfo, err := net.Interfaces()
	infos := make([]interface{}, 0)
	if err == nil && len(interfaceInfo) != 0 {
		for _, info := range interfaceInfo {
			ipAddr4 := ""
			ipAddr6 := ""
			netmask := ""
			for _, addr := range info.Addrs {
				if strings.Contains(addr.Addr, ":") {
					ipAddr6 = addr.Addr
				}
				if strings.Contains(addr.Addr, ".") {
					ip, ipv4Net, err := gonet.ParseCIDR(addr.Addr)
					if err != nil {
						log.Fatal(err)
					}
					ipAddr4 = fmt.Sprintf("%s", ip)
					netmask = ipv4MaskString(ipv4Net.Mask)
				}
			}
			speed, err := getInterfaceSpeed(info.Name)
			if err != nil {
				log.Printf("Failed to get if %s speed.", info.Name)
			}
			ifMap := map[string]string{
				"name":         info.Name,
				"mtu":          strconv.Itoa(info.MTU),
				"hardwareAddr": info.HardwareAddr,
				"ipAddr4":      ipAddr4,
				"ipAddr6":      ipAddr6,
				"netmask":      netmask,
				"speed":        speed,
				"flags":        strings.Join(info.Flags, ","),
			}
			infos = append(infos, ifMap)
		}
		info["interfaceInfo"] = infos
	}
	return info
}

func getMemInfo() map[string]interface{} {
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info["memInfo"] = memInfo
	}
	swapInfo, err := mem.SwapMemory()
	if err == nil {
		info["swapInfo"] = swapInfo
	}
	return info
}

func getDiskInfo() map[string]interface{} {
	diskInfos, err := disk.Partitions(true)
	infos := make([]interface{}, 0)
	if err == nil && len(diskInfos) != 0 {
		for _, diskInfo := range diskInfos {
			usage, err := disk.Usage(diskInfo.Mountpoint)
			if err == nil && usage != nil {
				infoMap := map[string]string{
					"device": diskInfo.Device,
					"mountpoint": diskInfo.Mountpoint,
					"fstype": diskInfo.Fstype,
					"total": strconv.FormatUint(usage.Total, 10),
					"free": strconv.FormatUint(usage.Free, 10),
					"used": strconv.FormatUint(usage.Used, 10),
					"usedPercent": strconv.FormatFloat(usage.UsedPercent, 'f', -1, 64),
					"inodesTotal": strconv.FormatUint(usage.InodesTotal, 10),
					"inodesUsed": strconv.FormatUint(usage.InodesUsed, 10),
					"inodesFree": strconv.FormatUint(usage.InodesFree, 10),
					"inodesUsedPercent": strconv.FormatFloat(usage.InodesUsedPercent, 'f', -1, 64),
				}
				infos = append(infos, infoMap)
			}
		}
	}
	info["diskInfo"] = infos
	return info
}

func getProcessInfo(proc *process.Process) map[string]interface{} {
	processInfo := make(map[string]interface{})
	processInfo["pid"] = proc.Pid
	processInfo["name"], _ = proc.Name()
	processInfo["status"], _ = proc.Status()
	processInfo["cpuPercent"], _ = proc.CPUPercent()
	if times, err := proc.Times(); err == nil {
		if times != nil {
			processInfo["cpuTimes"] = times.Total()
		}
	}
	processInfo["memoryPercent"], _ = proc.MemoryPercent()
	processInfo["create_time"], _ = proc.CreateTime()
	processInfo["workspace"], _ = proc.Cwd()
	processInfo["execPath"], _ = proc.Exe()
	processInfo["owner"], _ = proc.Username()
	processInfo["cmdLine"], _ = proc.Cmdline()
	return processInfo
}

func getProcsInfo() map[string]interface{} {
	var procs []int32
	var err1 error
	procs, err1 = process.Pids()
	if err1 != nil || len(procs) <= 0 {
		return info
	}
	procInfo := make(map[string]interface{})
	for _, pid := range procs {
		proc, _ := process.NewProcess(int32(pid))
		procInfo[strconv.Itoa(int(pid))] = getProcessInfo(proc)
	}
	info["procInfo"] = procInfo
	return info
}

func GetAllInfo() string {
	info = make(map[string]interface{})
	info["uuid"] = g.Config().Hostname
	info["tname"] = "Host"
	info["timeout"] = g.Config().Consul.Timeout
	info["timestamp"] = time.Now().Unix()
	getHostInfo()
	getLoadInfo()
	getCpuInfo()
	getMemInfo()
	getInterfaceInfo()
	getDiskInfo()
	//getProcsInfo()

	var data []byte
	data, err := json.Marshal(info)
	if err == nil {
		return fmt.Sprintf("{\"auto\":%s}", string(data))
	} else {
		return ""
	}
}

func sendSystemInfoToConsul(consulUrl string, nodeId string, data string) {
	if data != "" {
		url := fmt.Sprintf("%s/v1/kv/object/%s/system?raw", consulUrl, nodeId)
		r := httplib.Put(url)
		r.Body(data)
		ret, err := r.String()
		fmt.Println(ret, err)
	}
}

func formatDeltatime(sec uint64) (string, error){
	if sec < 0 {
		err := errors.New("The deltatime must be positive.")
		return "", err
	}
	days := sec / (60 * 60 * 24)
	hours := (sec % (60 * 60 * 24)) / (60 * 60)
	mins := (sec % (60 * 60)) / 60
	secs := sec % 60
	formatTime := fmt.Sprintf("%d days %d:%d:%d", days, hours, mins, secs)
	return formatTime, nil
}

func getInterfaceSpeed(inf string) (string, error) {
	ifSpeedFile := fmt.Sprintf("/sys/class/net/%s/speed", inf)

	data, err := ioutil.ReadFile(ifSpeedFile)
	if err != nil {
		fmt.Println(err)
		return "Not support!", err
	}
	//speed := string(data)
	var speed string
	s := strings.TrimSpace(string(data))
	if s == "-1" {
		speed = "Unknown"
	} else {
		speed = s + "Mb/s"
	}
	fmt.Println(speed)
	return speed, err
}