package util

import (
	libvirt "github.com/libvirt/libvirt-go"
)

func GetHostInfo(domain *libvirt.Domain) (map[string]interface{}, error) {
	info := make(map[string]interface{})
	hostname, err := domain.GetHostname(0)
	if err == nil && hostname != "" {
		info["hostname"] = hostname
	}

	domainname, _ := domain.GetName()
	if err == nil && domainname != "" {
		info["domainname"] = domainname
	}
	return info, nil
}
