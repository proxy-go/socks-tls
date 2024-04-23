package auth

import (
	"testing"
)

func TestLdapVerify(t *testing.T) {
	username := "test"
	password := "password"
	flag := LdapVerify("192.168.1.211:3890", "dc=windvpn,dc=com", username, password)
	t.Log(flag)
}
