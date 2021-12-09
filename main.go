package main

import (
	"github.com/leki75/go-tcp/producer"
	"github.com/leki75/go-tcp/tcp"
)

func main() {
	ch := make(chan []byte, 10000)
	go producer.TradeProducer(ch)
	server, err := tcp.NewServer(ch)
	if err != nil {
		panic(err)
	}
	err = server.Listen(":8080")
	panic(err)
}
