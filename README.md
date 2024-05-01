# socks-tls

A simple socks5 proxy over tls

# Features
* Support socks5(tcp/udp)
* Support socks5 over tls
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
  -tls 
        enable tls
  -auto
        tls cert auto mode
  -cert string
        certificate file (default "")
  -key string
        private key file (default "")
  -d string
        tls domain (default "localhost")
  -t int
        dial timeout in seconds (default 30)
  -http string
        http server address (default ":80")
  -ldap
        enable ldap auth
  -ldap-addr string
        ldap address (default "127.0.0.1:3890")
  -ldap-base-dn string
        ldap base dn (default "dc=example,dc=com")
  -iface string
        specified interface
```

# Run socks-tls with docker

## no auth
```
docker run  -d --restart=always --net=host \
-p 1080:1080 -p 1080:1080/udp --name socks-tls proxygo/socks-tls -l :1080
```

## auth
```
docker run  -d --restart=always --net=host \
-p 1080:1080 -p 1080:1080/udp --name socks-tls proxygo/socks-tls -l :1080 -u root -p 123456
```

## tls auth
```
docker run  -d --restart=always --net=host \
-p 1080:1080 -p 1080:1080/udp --name socks-tls proxygo/socks-tls -l :1080 -u root -p 123456 -tls
```

## automatic tls using Let's Encrypt
docker run  -d --restart=always --net=host \
-p 1080:1080 -p 1080:1080/udp --name socks-tls proxygo/socks-tls -l :1080 -u root -p 123456 -tls -auto
```

## specified interface
```
docker run  -d --restart=always --net=host \
-p 1080:1080 -p 1080:1080/udp --name socks-tls proxygo/socks-tls -l :1080 -iface tun0
```

## ldap auth
```
docker run  -d --restart=always --net=host \
-p 1080:1080 -p 1080:1080/udp --name socks-tls proxygo/socks-tls -l :1080  -ldap -ldap-addr 127.0.0.1:3890 -ldap-base-dn dc=example,dc=com
```

