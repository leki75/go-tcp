package producer

import (
	"testing"

	"github.com/leki75/go-tcp/config"
	"github.com/leki75/go-tcp/proto"
)

func BenchmarkTCPTradeProducer_binary(b *testing.B) {
	config.Encoding = config.EncodingBinary
	ch := make(chan []byte, 1)
	go TCPTradeProducer(ch)
	for i := 0; i < b.N; i++ {
		<-ch
	}
}

func BenchmarkTCPTradeProducer_text(b *testing.B) {
	config.Encoding = config.EncodingText
	ch := make(chan []byte, 1)
	go TCPTradeProducer(ch)
	for i := 0; i < b.N; i++ {
		<-ch
	}
}

func BenchmarkGRPCTradeProducer(b *testing.B) {
	ch := make(chan *proto.Trade, 1)
	go GRPCTradeProducer(ch)
	for i := 0; i < b.N; i++ {
		<-ch
	}
}
