package binary

import (
	"encoding/binary"
	"math"

	"github.com/leki75/go-tcp/schema/raw"
)

func MarshalTrade(t *raw.Trade) []byte {
	b := make([]byte, 54)
	b[0] = 't' // type
	copy(b[1:], t.Symbol[:])
	copy(b[12:], t.Conditions[:])
	binary.BigEndian.PutUint64(b[16:], t.Id)
	binary.BigEndian.PutUint64(b[24:], uint64(t.Timestamp))
	binary.BigEndian.PutUint64(b[32:], uint64(t.ReceivedAt))
	binary.BigEndian.PutUint64(b[40:], math.Float64bits(t.Price))
	binary.BigEndian.PutUint32(b[48:], t.Size)
	b[52] = t.Exchange
	b[53] = t.Tape
	return b
}

func UnmarshalTrade(b []byte) *raw.Trade {
	t := &raw.Trade{
		Id:         binary.BigEndian.Uint64(b[16:]),
		Timestamp:  int64(binary.BigEndian.Uint64(b[24:])),
		ReceivedAt: int64(binary.BigEndian.Uint64(b[32:])),
		Price:      math.Float64frombits(binary.BigEndian.Uint64(b[40:])),
		Size:       binary.BigEndian.Uint32(b[48:]),
		Exchange:   b[52],
		Tape:       b[53],
	}
	copy(t.Symbol[:], b[1:12])
	copy(t.Conditions[:], b[12:16])
	return t
}
