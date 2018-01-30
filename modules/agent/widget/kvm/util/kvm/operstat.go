package util

import (
	libvirt "github.com/libvirt/libvirt-go"
	"log"
)

func GetOperStatus(domain *libvirt.Domain) (string, error) {
	var operStatus string = "unknown"
	domainInfo, err := domain.GetInfo()
	if err != nil || domainInfo == nil {
		log.Println("get domain info error", err)
		return operStatus, err
	}

	switch domainInfo.State {
	case libvirt.DOMAIN_NOSTATE:
		operStatus = "unknown"
	case libvirt.DOMAIN_RUNNING:
		operStatus = "running"
	case libvirt.DOMAIN_BLOCKED:
		operStatus = "blocked"
	case libvirt.DOMAIN_PAUSED:
		operStatus = "paused"
	case libvirt.DOMAIN_SHUTDOWN:
		operStatus = "shutdown"
	case libvirt.DOMAIN_CRASHED:
		operStatus = "crashed"
	case libvirt.DOMAIN_PMSUSPENDED:
		operStatus = "pmsuspended"
	case libvirt.DOMAIN_SHUTOFF:
		operStatus = "shutoff"
	}
	return operStatus, nil
}
