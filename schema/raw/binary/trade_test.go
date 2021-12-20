package binary

import (
	"encoding/binary"
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

func BenchmarkUnmarshalTrade(b *testing.B) {
	trade := raw.Trade{
		Id:         0,
		Timestamp:  time.Now().UnixNano(),
		ReceivedAt: time.Now().UnixNano(),
		Symbol:     [11]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1'},
		Price:      123.456,
		Size:       100,
		Conditions: [4]byte{'1', '2', '3', '4'},
		Exchange:   'N',
		Tape:       'C',
	}
	buf := MarshalTrade(&trade)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint64(buf[17:], uint64(i))
		_ = UnmarshalTrade(buf)
	}

}

func TestMarshalUnmarshalTrade(t *testing.T) {
	now := time.Now().UnixNano()
	expected := &raw.Trade{
		Id:         1234567890,
		Price:      123.456,
		Timestamp:  now - int64(time.Millisecond),
		ReceivedAt: now,
		Size:       987,
		Symbol:     [11]byte{'A', 'A', 'P', 'L'},
		Conditions: [4]byte{'@'},
		Exchange:   'X',
		Tape:       'T',
	}
	b := MarshalTrade(expected)
	got := UnmarshalTrade(b)
	assert.EqualValues(t, expected, got)
}
