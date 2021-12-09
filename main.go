package main

import (
	"github.com/leki75/go-tcp/producer"
	"github.com/leki75/go-tcp/tcp"
)

func main() {
	ch := make(chan []byte, 10000)
	go tcp.NewServer(ch).Listen(":8080")
	producer.TradeProducer(ch)
}
