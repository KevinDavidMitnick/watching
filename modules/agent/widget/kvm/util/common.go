package util

import (
	"bytes"
	"encoding/json"
	//	"flag"
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"io/ioutil"
	"log"
	"net/http"
)

func ConnectLibvirt(addr string) (*libvirt.Connect, error) {
	conn, err := libvirt.NewConnect("qemu+tcp://root@" + addr + "/system")
	if err != nil {
		log.Fatalln("Can not connect libvirt,err is:", err)
	}
	return conn, err
}

func ListAliveDomains(conn *libvirt.Connect) ([]libvirt.Domain, error) {
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	return doms, err
}

func ListAllDomains(conn *libvirt.Connect) ([]libvirt.Domain, error) {
	doms, err := conn.ListAllDomains(0)
	return doms, err
}

func GetValueFromQga(uuid string, action string) (interface{}, error) {
	//cfg := flag.String("c", "/home/shadowalker/go/src/gitlab.chinac.com/OpsUltra/tiresias-probe/kvmcollector/cfg.json", "configuration file")
	//flag.Parse()
	//g.ParseConfig(*cfg)

	host := g.Config().Kvm.Host
	port := g.Config().Kvm.Port
	timeout := g.Config().Kvm.Timeout

	url := fmt.Sprintf("http://%s:%d", host, port)
	bodyStr := fmt.Sprintf("{\"%s\": [\"{\\\"execute\\\": \\\"%s\\\", \\\"arguments\\\": {\\\"id\\\": \\\"%s\\\"}}\"]}", uuid, action, uuid)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyStr)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("TIMEOUT", fmt.Sprintf("%d", timeout))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret []map[string]map[string]interface{}
	json.Unmarshal(body, &ret)

	if len(ret) < 1 {
		return nil, nil
	}

	retStatus, stat := ret[0]["return"]
	if stat == false {
		return nil, nil
	}
	msg, _ := retStatus["msg"]
	if msg != "Sucess" || msg == "timeout" {
		return nil, nil
	}

	retValue, stat := retStatus["value"]
	if stat == false {
		return nil, nil
	}
	return retValue, nil
}
