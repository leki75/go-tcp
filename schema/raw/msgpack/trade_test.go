package msgpack

import (
	"testing"
	"time"

	"github.com/leki75/go-tcp/schema/raw"
)

func BenchmarkMarshalTrade(b *testing.B) {
	trade := raw.Trade{
		Symbol:     [11]byte{'A', 'A', 'P', 'L'},
		Price:      123.456,
		Size:       100,
		Conditions: [4]byte{'@'},
		Exchange:   'N',
		Tape:       'C',
	}
	for i := 0; i < b.N; i++ {
		trade.Id = uint64(i)
		trade.Timestamp = time.Now().UnixNano()
		trade.ReceivedAt = trade.Timestamp
		_ = MarshalTrade(&trade)
	}
}
