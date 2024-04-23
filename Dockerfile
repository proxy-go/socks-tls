FROM golang:alpine

WORKDIR /app
COPY . /app
ENV GO111MODULE=on
RUN go build -o ./bin/socks-tls ./main.go

ENTRYPOINT ["./bin/socks-tls"]

