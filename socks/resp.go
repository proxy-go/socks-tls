package socks

import (
	"bytes"
	"encoding/binary"
	"net"
)

/*
*

	+----+--------+
	|VER | METHOD |
	+----+--------+
	| 1  |   1    |
	+----+--------+
*/
func respAuthType(conn net.Conn, authType uint8) {
	conn.Write([]byte{Socks5Version, authType})
}

/*
*

	+----+--------+
	|VER | STATUS |
	+----+--------+
	| 1  |   1    |
	+----+--------+
*/
func respAuthStatus(conn net.Conn, status uint8) {
	conn.Write([]byte{UserAuthVersion, status})
}

/*
*

	+----+-----+-------+------+----------+----------+
	|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/
func resp(conn net.Conn, rep byte) {
	conn.Write([]byte{Socks5Version, rep, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

/*
*

	+----+-----+-------+------+----------+----------+
	|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/
func respUDP(conn net.Conn, bindAddr *net.UDPAddr) {
	resp := []byte{Socks5Version, SuccessReply, 0x00, 0x01}
	buffer := bytes.NewBuffer(resp)
	binary.Write(buffer, binary.BigEndian, bindAddr.IP.To4())
	binary.Write(buffer, binary.BigEndian, uint16(bindAddr.Port))
	conn.Write(buffer.Bytes())
}
