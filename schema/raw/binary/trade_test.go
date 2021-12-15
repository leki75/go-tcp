package binary

import (
	"testing"
	"time"

	"github.com/leki75/go-tcp/schema/raw"
	"github.com/stretchr/testify/assert"
)

func BenchmarkMarshalTrade(b *testing.B) {
	trade := raw.Trade{
		Symbol:     [11]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1'},
		Price:      123.456,
		Size:       100,
		Conditions: [4]byte{'1', '2', '3', '4'},
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

func TestMarshalUnmarshalTrade(t *testing.T) {
	now := time.Now().UnixNano()
	expected := &raw.Trade{
		Price:      123.456,
		Timestamp:  now - int64(time.Millisecond),
		ReceivedAt: now,
		Size:       987,
		Symbol:     *(*[11]byte)([]byte("12345678901")),
		Conditions: [4]byte{'1', '2', '3', '4'},
		Exchange:   'X',
		Tape:       'T',
	}
	b := MarshalTrade(expected)
	got := UnmarshalTrade(b)
	assert.EqualValues(t, expected, got)
}
