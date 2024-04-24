package socks

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/proxy-go/socks-tls/auth"
	"github.com/proxy-go/socks-tls/certs"
)

// Tcp server struct
type TCPServer struct {
	config   Config
	udpConn  *net.UDPConn
	publicIP string
	outIP    net.IP
	outIface *net.Interface
}

// Start tcp server
func (t *TCPServer) Start() {
	var l net.Listener
	var err error
	var cert tls.Certificate
	if t.config.TLS {
		if t.config.TLSCertFile != "" && t.config.TLSKeyFile != "" {
			cert, err = tls.LoadX509KeyPair(t.config.TLSCertFile, t.config.TLSKeyFile)
			if err != nil {
				log.Panic(err)
			}
		} else {
			cert = certs.GenerateCert("localhost")
		}

		c := &tls.Config{Certificates: []tls.Certificate{cert}}
		l, err = tls.Listen("tcp", t.config.LocalAddr, c)
		if err != nil {
			log.Panicf("[tls] failed to listen tcp %v", err)
		}
		log.Printf("socks-tls tls proxy started on %s", t.config.LocalAddr)
	} else {
		l, err = net.Listen("tcp", t.config.LocalAddr)
		if err != nil {
			log.Panicf("[tcp] failed to listen tcp %v", err)
		}
		log.Printf("socks-tls tcp proxy started on %s", t.config.LocalAddr)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go t.handler(conn)
	}
}

// Tcp handler
func (t *TCPServer) handler(conn net.Conn) {
	buf := make([]byte, BufferSize)
	// read version
	n, err := conn.Read(buf[0:])
	if err != nil || err == io.EOF || n == 0 {
		conn.Close()
		return
	}
	b := buf[0:n]
	if b[0] != Socks5Version {
		conn.Close()
		return
	}
	if !t.config.LdapAuth && t.config.Username == "" && t.config.Password == "" {
		// no auth
		respAuthType(conn, NoAuth)
	} else {
		// username and password auth
		respAuthType(conn, UserPassAuth)
		username, password := t.getUserPwd(conn)
		if t.config.LdapAuth {
			ok := auth.LdapVerify(t.config.LdapAddr, t.config.LdapBaseDN, username, password)
			if ok {
				respAuthStatus(conn, AuthSuccess)
			} else {
				respAuthStatus(conn, AuthFailure)
			}
		} else {
			if username == t.config.Username && password == t.config.Password {
				respAuthStatus(conn, AuthSuccess)
			} else {
				respAuthStatus(conn, AuthFailure)
			}
		}
	}
	// read cmd
	n, err = conn.Read(buf[0:])
	if err != nil || err == io.EOF || n < 2 {
		conn.Close()
		return
	}
	b = buf[0:n]
	switch b[1] {
	case ConnectCommand:
		t.TCPProxy(conn, b)
	case AssociateCommand:
		t.UDPProxy(conn, t.udpConn, t.config)
	case BindCommand:
		resp(conn, CommandNotSupported)
	default:
		resp(conn, CommandNotSupported)
	}
}

/*
*
  - Get username and password from conn
    +----+------+----------+------+----------+
    |VER | ULEN | UNAME | PLEN | PASSWD |
    +----+------+----------+------+----------+
    | 1 | 1 | 1 to 255 | 1 | 1 to 255 |
    +----+------+----------+------+----------+
*/
func (t *TCPServer) getUserPwd(conn net.Conn) (user, pwd string) {
	ver := make([]byte, 1)
	n, err := conn.Read(ver)
	if err != nil || n == 0 {
		return "", ""
	}
	if uint(ver[0]) != uint(UserAuthVersion) {
		return "", ""
	}
	ulen := make([]byte, 1)
	n, err = conn.Read(ulen)
	if err != nil || n == 0 {
		return "", ""
	}
	if uint(ulen[0]) < 1 {
		return "", ""
	}
	uname := make([]byte, uint(ulen[0]))
	n, err = conn.Read(uname)
	if err != nil || n == 0 {
		return "", ""
	}
	user = string(uname)

	plen := make([]byte, 1)
	n, err = conn.Read(plen)
	if err != nil || n == 0 {
		return "", ""
	}
	if uint(plen[0]) < 1 {
		return "", ""
	}
	passwd := make([]byte, uint(plen[0]))
	n, err = conn.Read(passwd)
	if err != nil || n == 0 {
		return "", ""
	}
	pwd = string(passwd)
	return user, pwd
}

/*
*
  - Get host and port from data
    +----+-----+-------+------+----------+----------+
    |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
    +----+-----+-------+------+----------+----------+
    | 1  |  1  | X'00' |  1   | Variable |    2     |
    +----+-----+-------+------+----------+----------+
*/
func (t *TCPServer) getAddr(b []byte) (host string, port string) {
	len := len(b)
	if len < 4 {
		return "", ""
	}
	switch b[3] {
	case Ipv4Address:
		host = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case FqdnAddress:
		host = string(b[5 : len-2])
	case Ipv6Address:
		host = net.IP(b[4:20]).String()
	default:
		return "", ""
	}
	port = strconv.Itoa(int(b[len-2])<<8 | int(b[len-1]))
	return host, port
}

// Tcp proxy
func (t *TCPServer) TCPProxy(conn net.Conn, data []byte) {
	host, port := t.getAddr(data)
	if host == "" || port == "" {
		conn.Close()
		return
	}
	remoteConn, err := dial("tcp", net.JoinHostPort(host, port), t.outIface, t.outIP, t.config.Timeout)
	if err != nil {
		log.Printf("[tcp] failed to dial tcp %v", err)
		resp(conn, ConnectionRefused)
		return
	}
	// resp tcp connect success
	resp(conn, SuccessReply)
	go copy(remoteConn, conn)
	copy(conn, remoteConn)
}

// Udp proxy
func (t *TCPServer) UDPProxy(tcpConn net.Conn, udpConn *net.UDPConn, config Config) {
	defer tcpConn.Close()
	if udpConn == nil {
		log.Printf("[udp] failed to start udp server on %v", config.LocalAddr)
		return
	}
	bindAddr, _ := net.ResolveUDPAddr("udp", udpConn.LocalAddr().String())
	if bindAddr.IP.To4() == nil {
		if t.publicIP == "" {
			t.publicIP = getPublicIP()
		}
		bindAddr.IP = net.ParseIP(t.publicIP)
	}
	// resp udp associate
	respUDP(tcpConn, bindAddr)
	// keep tcp conn alive
	done := make(chan bool)
	if config.TLS {
		go keepTLSAlive(tcpConn.(*tls.Conn), done)
	} else {
		go keepTCPAlive(tcpConn.(*net.TCPConn), done)
	}
	<-done
}

// Get public ip
func getPublicIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

// Keep tcp connection alive
func keepTCPAlive(tcpConn *net.TCPConn, done chan<- bool) {
	tcpConn.SetKeepAlive(true)
	buf := make([]byte, BufferSize)
	for {
		_, err := tcpConn.Read(buf[0:])
		if err != nil {
			break
		}
	}
	done <- true
}

// Keep tls connection alive
func keepTLSAlive(conn *tls.Conn, done chan<- bool) {
	buf := make([]byte, BufferSize)
	for {
		_, err := conn.Read(buf[0:])
		if err != nil {
			break
		}
	}
	done <- true
}

// Copy data from src to dst
func copy(to io.WriteCloser, from io.ReadCloser) {
	defer to.Close()
	defer from.Close()
	io.Copy(to, from)
}
