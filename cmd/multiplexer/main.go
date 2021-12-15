package main

import (
	"github.com/leki75/go-tcp/config"
	grpcclient "github.com/leki75/go-tcp/grpc/client"
	"github.com/leki75/go-tcp/schema/proto"
	tcpclient "github.com/leki75/go-tcp/tcp/client"
	tcpserver "github.com/leki75/go-tcp/tcp/server"
)

func main() {
	var err error

	switch config.Proto {
	case config.ProtoTCP:
		ch := make(chan []byte, 10000)
		go tcpserver.NewServer(ch).Listen(":8000")
		err = tcpclient.NewClient(ch).Connect(":8080")

	case config.ProtoGRPC:
		ch := make(chan *proto.Trade, 10000)
		go tcpserver.NewProtoServer(ch).Listen(":8000")
		err = grpcclient.NewClient(ch).Connect(":9090")
	}

	panic(err)
}
