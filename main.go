package main

import (
	"flag"

	"github.com/proxy-go/socks-tls/socks"
)

func main() {
	config := socks.Config{}
	flag.StringVar(&config.LocalAddr, "l", ":1080", "local address")
	flag.StringVar(&config.Username, "u", "", "username")
	flag.StringVar(&config.Password, "p", "", "password")
	flag.StringVar(&config.TLSCertFile, "cert", "./certs/certificate.crt", "certificate file")
	flag.StringVar(&config.TLSKeyFile, "key", "./certs/private.key", "private key file")
	flag.BoolVar(&config.TLS, "tls", false, "enable tls")
	flag.StringVar(&config.Iface, "iface", "", "specified interface")
	flag.IntVar(&config.Timeout, "t", 30, "dial timeout in seconds")
	flag.BoolVar(&config.LdapAuth, "ldap", false, "enable ldap auth")
	flag.StringVar(&config.LdapAddr, "ldap-addr", "127.0.0.1:3890", "ldap address")
	flag.StringVar(&config.LdapBaseDN, "ldap-base-dn", "dc=example,dc=com", "ldap base dn")
	flag.Parse()

	socks.Start(config)
}
