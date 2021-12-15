package main

import (
	"log"

	"github.com/leki75/go-tcp/config"
	grpcserver "github.com/leki75/go-tcp/grpc/server"
	"github.com/leki75/go-tcp/producer"
	"github.com/leki75/go-tcp/schema/proto"
	tcpserver "github.com/leki75/go-tcp/tcp/server"
)

func main() {
	var err error
	switch config.Proto {
	case config.ProtoTCP:
		log.Println("Proto TCP")
		ch := make(chan []byte, 10000)
		go producer.ByteSliceTradeProducer(ch)
		err = tcpserver.NewServer(ch).Listen(":8080")

	case config.ProtoGRPC:
		log.Println("Proto gRPC")
		ch := make(chan *proto.Trade, 10000)
		go producer.ProtobufTradeProducer(ch)
		err = grpcserver.NewServer(ch).Listen(":9090")
	}

	panic(err)
}
