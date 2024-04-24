package socks

func Start(config Config) {
	outIP, outIface := GetIface(config)

	u := &UDPRelay{config: config, outIP: outIP, outIface: outIface}
	udpConn := u.Start()

	t := &Socks5Server{config: config, udpConn: udpConn, outIP: outIP, outIface: outIface}
	t.Start()
}
