# socks-tls

A socks5 proxy over tls

# Features
* Support connect
* Support udp associate
* Support tcp over tls
* Support specified interface
* Support ldap auth

# Usage
```
Usage of /main:
  -l string
        local address (default ":1080")
  -p string
        password
  -u string
        username
  -sk string
        server key file path (default "../certs/server.key")
  -sp string
        server pem file path (default "../certs/server.pem")
  -tls enable tls
  -iface string
        specified interface
  -t int
        dial timeout in seconds (default 30)
  -ldap
        enable ldap auth
  -ldap-addr string
        ldap address (default "127.0.0.1:3890")
  -ldap-base-dn string
        ldap base dn (default "dc=example,dc=com")
```



# License
[The MIT License (MIT)](https://raw.githubusercontent.com/proxy-go/socks-tls/main/LICENSE)
