package json

import (
	"testing"
	"time"

	"github.com/leki75/go-tcp/schema/proto"
)

func BenchmarkMarshalTradeNew(b *testing.B) {
	trade := proto.Trade{
		Symbol:     "12345678901",
		Price:      123.456,
		Volume:     100,
		Conditions: []byte{'1', '2', '3', '4'},
		Exchange:   'N',
		Tape:       'C',
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trade.Id = uint64(i)
		trade.ReceivedAt = time.Now().UnixNano()
		trade.Timestamp = uint64(trade.ReceivedAt)
		_ = MarshalTrade(&trade)
	}
}
