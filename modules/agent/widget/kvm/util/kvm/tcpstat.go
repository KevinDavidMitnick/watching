package util

import (
	"encoding/json"
	libvirt "github.com/libvirt/libvirt-go"
)

type TcpStat struct {
	Name  string  `json:"name"`
	Count float64 `json:"count"`
}

func GetTcpStat(domain *libvirt.Domain) ([]int64, error) {
	var tcps []int64
	uuid, err := domain.GetUUIDString()
	if err != nil {
		return nil, err
	}
	tcpStat, err := GetValueFromQga(uuid, "guest-get-tcp-number")
	if tcpStat == nil || err != nil {
		return nil, err
	}
	var stat []TcpStat
	data, err := json.Marshal(tcpStat)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &stat)
	var established, synrecv, timewait, closewait, listen int64
	for _, tcp := range stat {
		switch tcp.Name {
		case "ESTABLISHED":
			established = int64(tcp.Count)
		case "SYN_RECV":
			synrecv = int64(tcp.Count)
		case "TIME_WAIT":
			timewait = int64(tcp.Count)
		case "CLOSE_WAIT":
			closewait = int64(tcp.Count)
		case "LISTEN":
			listen = int64(tcp.Count)
		}
		tcps = []int64{established, synrecv, timewait, closewait, listen}
	}
	return tcps, nil
}
