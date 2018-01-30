package util

import (
	"github.com/beevik/etree"
	libvirt "github.com/libvirt/libvirt-go"
	"time"
)

func GetInterfaceIO(domain *libvirt.Domain) (map[string][]int64, error) {
	devs := make(map[string][]int64)
	xml, err := domain.GetXMLDesc(0)
	if err != nil {
		return devs, err
	}

	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return devs, err
	}

	for _, elem := range doc.FindElements("//devices/interface/target") {
		for _, attr := range elem.Attr {
			if attr.Key == "dev" {
				interfaceStatsBegin, err := domain.InterfaceStats(attr.Value)
				if err != nil {
					return devs, err
				}
				time.Sleep(time.Duration(1) * time.Second)
				interfaceStatsEnd, err := domain.InterfaceStats(attr.Value)
				if err != nil {
					return devs, err
				}
				rxBytes := int64(interfaceStatsEnd.RxBytes - interfaceStatsBegin.RxBytes)
				rxPackets := int64(interfaceStatsEnd.RxPackets - interfaceStatsBegin.RxPackets)
				txBytes := int64(interfaceStatsEnd.TxBytes - interfaceStatsBegin.TxBytes)
				txPackets := int64(interfaceStatsEnd.TxPackets - interfaceStatsBegin.TxPackets)
				devs["iface="+attr.Value] = []int64{rxBytes, rxPackets, txBytes, txPackets}
			}
		}
	}
	return devs, nil
}

func GetInterfaceInfo(domain *libvirt.Domain) map[string]interface{} {
	M := make(map[string]interface{})
	devs := make(map[string]interface{})
	xml, err := domain.GetXMLDesc(0)
	if err != nil {
		return devs
	}

	doc := etree.NewDocument()
	err = doc.ReadFromString(xml)
	if err != nil {
		return devs
	}

	root := doc.SelectElement("domain")

	if root == nil {
		return devs
	}

	root = root.SelectElement("devices")
	if root == nil {
		return devs
	}

	for _, inter := range root.SelectElements("interface") {
		if target := inter.SelectElement("target"); target != nil {
			if mac := inter.SelectElement("mac"); mac != nil {
				dev := target.SelectAttrValue("dev", "")
				addr := mac.SelectAttrValue("address", "")
				if dev != "" && addr != "" {
					devs[dev] = map[string]interface{}{"hardwareAddr": addr}
				}
			}
		}
	}
	M["interfaceInfo"] = devs
	return M
}
