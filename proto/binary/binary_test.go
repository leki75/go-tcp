package binary

import (
	"testing"
	"time"

	"github.com/leki75/go-tcp/proto"
	"github.com/stretchr/testify/assert"
)

func BenchmarkMarshalTrade(b *testing.B) {
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

func TestMarshalUnmarshalTrade(t *testing.T) {
	now := time.Now().UnixNano()
	expected := &proto.Trade{
		Price:      123.456,
		Timestamp:  uint64(now) - uint64(time.Millisecond),
		ReceivedAt: now,
		Volume:     987,
		Symbol:     "AAPL",
		Conditions: []byte{'@'},
		Exchange:   'X',
		Tape:       'T',
	}
	b := MarshalTrade(expected)
	got := UnmarshalTrade(b)
	assert.EqualValues(t, expected, got)
}

func TestGetConditions(t *testing.T) {
	conditions := [4]byte{}
	assert.Equal(t, []byte{0}, GetConditions(conditions[:]))

	conditions = [4]byte{'A'}
	assert.Equal(t, []byte{'A'}, GetConditions(conditions[:]))

	conditions = [4]byte{0, 'A', 0, 0}
	assert.Equal(t, []byte{0, 'A'}, GetConditions(conditions[:]))

	conditions = [4]byte{'@', 'A', 'B', 'C'}
	assert.Equal(t, []byte{'@', 'A', 'B', 'C'}, GetConditions(conditions[:]))
}

func TestGetSymbol(t *testing.T) {
	symbol := [11]byte{}
	assert.Equal(t, "", GetSymbol(symbol[:]))

	symbol = [11]byte{'A', 'A', 'P', 'L'}
	assert.Equal(t, "AAPL", GetSymbol(symbol[:]))

	symbol = [11]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1'}
	assert.Equal(t, "12345678901", GetSymbol(symbol[:]))
}
