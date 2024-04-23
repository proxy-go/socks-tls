package socks

import (
	"log"
	"net"
)

func Start(config Config) {
	outIP, outIface := getIface(config)
	// start udp server
	u := &UDPServer{config: config, outIP: outIP, outIface: outIface}
	udpConn := u.Start()
	// start tcp server
	t := &TCPServer{config: config, udpConn: udpConn, outIP: outIP, outIface: outIface}
	t.Start()
}

func getIface(config Config) (net.IP, *net.Interface) {
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
