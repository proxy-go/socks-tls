package socks

type Config struct {
	LocalAddr   string
	Username    string
	Password    string
	TLSKeyFile  string
	TLSCertFile string
	TLSDomain   string
	TLS         bool
	TLSAuto     bool
	Iface       string
	HttpAddr    string
	Timeout     int
	LdapAuth    bool
	LdapAddr    string
	LdapBaseDN  string
}
