package socks

type Config struct {
	LocalAddr   string
	Username    string
	Password    string
	TLSKeyFile  string
	TLSCertFile string
	TLS         bool
	Iface       string
	Timeout     int
	LdapAuth    bool
	LdapAddr    string
	LdapBaseDN  string
}
