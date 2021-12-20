package producer

import (
	"testing"

	"github.com/leki75/go-tcp/config"
	"github.com/leki75/go-tcp/schema/proto"
)

func BenchmarkByteSliceRawTradeProducer(b *testing.B) {
	config.Encoding = config.Binary
	ch := make(chan []byte, 1)
	go ByteSliceRawTradeProducer(ch)
	for i := 0; i < b.N; i++ {
		<-ch
	}
}

func BenchmarkProtobufTradeProducer(b *testing.B) {
	ch := make(chan *proto.Trade, 1)
	go ProtobufTradeProducer(ch)
	for i := 0; i < b.N; i++ {
		<-ch
	}
}
