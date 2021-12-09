package producer

import (
	"testing"
	"time"
)

func BenchmarkTrade_MarshalBinary(b *testing.B) {
	trade := Trade{
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

		_, err := trade.MarshalBinary()
		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkTrade_MarshalJSON(b *testing.B) {
	trade := Trade{
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

		_, err := trade.MarshalJSON()
		if err != nil {
			b.Fail()
		}
	}
}
