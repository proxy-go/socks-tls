#!bin/bash

#Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/socks-tls-linux-amd64 ./main.go
#Linux arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/socks-tls-linux-arm64 ./main.go
#Mac amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/socks-tls-darwin-amd64 ./main.go
#Mac arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/socks-tls-darwin-arm64 ./main.go
#Windows amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/socks-tls-windows-amd64.exe ./main.go
#Windows arm64
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./bin/socks-tls-windows-arm64.exe ./main.go
echo "DONE!!!"
