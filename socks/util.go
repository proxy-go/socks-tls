package socks

import (
	"log"
	"net"
)

func GetIface(config Config) (net.IP, *net.Interface) {
	if config.Iface == "" {
		return nil, nil
	}
	ief, err := net.InterfaceByName(config.Iface)
	if err != nil {
		log.Fatal("get Interface error", err)
	}
	addrs, err := ief.Addrs()
	if err != nil {
		log.Fatal(err)
	}
	outIface := ief
	outIP := addrs[0].(*net.IPNet).IP.To4()
	log.Printf("iface name: %v, out ip: %v", config.Iface, outIP.String())
	return outIP, outIface
}

func RecoverFromPanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic:%v", r)
	}
}
