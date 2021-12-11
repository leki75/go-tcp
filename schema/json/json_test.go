package json

import (
	"testing"
	"time"

	"github.com/leki75/go-tcp/schema"
)

func BenchmarkTrade_MarshalJSON(b *testing.B) {
	trade := schema.Trade{
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

		_, err := MarshalTrade(&trade)
		if err != nil {
			b.Fail()
		}
	}
}
