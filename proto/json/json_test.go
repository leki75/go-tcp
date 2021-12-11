package json

import (
	"testing"
	"time"

	"github.com/leki75/go-tcp/proto"
)

func BenchmarkMarshalTradeNew(b *testing.B) {
	trade := proto.Trade{
		Symbol:     "AAPL",
		Price:      123.456,
		Volume:     100,
		Conditions: []byte{'@'},
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
