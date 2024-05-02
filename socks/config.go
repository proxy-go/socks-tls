package socks

import (
	"fmt"
	"os"
)

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

func(config *Config) SetEnv(){
	localPort := os.Getenv("LOCAL_PORT")
	if localPort != "" {
		config.LocalAddr = fmt.Sprintf(":%v", localPort)
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort != "" {
		config.HttpAddr = fmt.Sprintf(":%v", httpPort)
	}
	username := os.Getenv("USERNAME")
	if username != "" {
		config.Username = username
	}
	password := os.Getenv("PASSWORD")
	if password != "" {
		config.Password = password
	}
	tls := os.Getenv("TLS")
	if tls == "true" {
		config.TLS = true
	}
	tlsAuto := os.Getenv("TLS_AUTO")
	if tlsAuto == "true" {
		config.TLSAuto = true
	}
	domain := os.Getenv("DOMAIN")
	if domain != "" {
		config.TLSDomain = domain
	}
}
