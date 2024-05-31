package socks

import (
	"bytes"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type UDPRelay struct {
	config      Config
	localConn   *net.UDPConn
	dstHeader   sync.Map
	remoteConns sync.Map
	outIP       net.IP
	outIface    *net.Interface
}

func (u *UDPRelay) Start() *net.UDPConn {
	udpAddr, _ := net.ResolveUDPAddr("udp", u.config.LocalAddr)
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("[udp] failed to listen udp %v", err)
		return nil
	}
	u.localConn = udpConn
	go u.toRemote()
	log.Printf("socks-tls udp relay started on %v", udpAddr)
	return u.localConn
}

func (u *UDPRelay) toRemote() {
	defer RecoverFromPanic()
	defer u.localConn.Close()
	buf := make([]byte, BufferSize)
	for {
		u.localConn.SetReadDeadline(time.Now().Add(time.Duration(u.config.Timeout) * time.Second))
		n, cliAddr, err := u.localConn.ReadFromUDP(buf)
		if err != nil || err == io.EOF || n == 0 {
			continue
		}
		b := buf[:n]
		dstAddr, header, data := u.getAddr(b)
		if dstAddr == nil || header == nil || data == nil {
			continue
		}
		key := cliAddr.String()
		value, ok := u.remoteConns.Load(key)
		if ok && value != nil {
			remoteConn := value.(*net.UDPConn)
			remoteConn.Write(data)
		} else {
			remoteConn, err := dial("udp", dstAddr.String(), u.outIface, u.outIP, u.config.Timeout)
			if remoteConn == nil || err != nil {
				log.Printf("failed to dial udp:%v", dstAddr)
				continue
			}
			u.remoteConns.Store(key, remoteConn)
			u.dstHeader.Store(key, header)
			go u.toLocal(remoteConn.(*net.UDPConn), cliAddr)
			remoteConn.Write(data)
		}
	}
}

func (u *UDPRelay) toLocal(remoteConn *net.UDPConn, cliAddr *net.UDPAddr) {
	defer remoteConn.Close()
	key := cliAddr.String()
	buf := make([]byte, BufferSize)
	remoteConn.SetReadDeadline(time.Now().Add(time.Duration(u.config.Timeout) * time.Second))
	for {
		n, _, err := remoteConn.ReadFromUDP(buf)
		if n == 0 || err != nil {
			break
		}
		if header, ok := u.dstHeader.Load(key); ok {
			var data bytes.Buffer
			data.Write(header.([]byte))
			data.Write(buf[:n])
			u.localConn.WriteToUDP(data.Bytes(), cliAddr)
		}
	}
	u.dstHeader.Delete(key)
	u.remoteConns.Delete(key)
}

/*
  - Get addr from packet
    +----+------+------+----------+----------+----------+
    |RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
    +----+------+------+----------+----------+----------+
    |  2 |   1  |   1  | Variable |     2    | Variable |
    +----+------+------+----------+----------+----------+
*/
func (u *UDPRelay) getAddr(b []byte) (dstAddr *net.UDPAddr, header []byte, data []byte) {
	if len(b) < 4 {
		return nil, nil, nil
	}
	if b[2] != 0x00 {
		log.Printf("[udp] not support frag %v", b[2])
		return nil, nil, nil
	}
	switch b[3] {
	case Ipv4Address:
		dstAddr = &net.UDPAddr{
			IP:   net.IPv4(b[4], b[5], b[6], b[7]),
			Port: int(b[8])<<8 | int(b[9]),
		}
		header = b[0:10]
		data = b[10:]
	case FqdnAddress:
		domainLength := int(b[4])
		domain := string(b[5 : 5+domainLength])
		ipAddr, err := net.ResolveIPAddr("ip", domain)
		if err != nil {
			log.Printf("[udp] failed to resolve dns %s:%v", domain, err)
			return nil, nil, nil
		}
		dstAddr = &net.UDPAddr{
			IP:   ipAddr.IP,
			Port: int(b[5+domainLength])<<8 | int(b[6+domainLength]),
		}
		header = b[0 : 7+domainLength]
		data = b[7+domainLength:]
	case Ipv6Address:
		{
			dstAddr = &net.UDPAddr{
				IP:   net.IP(b[4:20]),
				Port: int(b[20])<<8 | int(b[21]),
			}
			header = b[0:22]
			data = b[22:]
		}
	default:
		return nil, nil, nil
	}
	return dstAddr, header, data
}
