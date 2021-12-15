package binary

import (
	"encoding/binary"
	"math"

	"github.com/leki75/go-tcp/schema/raw"
)

func MarshalTrade(t *raw.Trade) []byte {
	b := make([]byte, 53)
	// TODO: add type
	binary.BigEndian.PutUint64(b[0:], t.Id)
	binary.BigEndian.PutUint64(b[8:], uint64(t.Timestamp))
	binary.BigEndian.PutUint64(b[16:], uint64(t.ReceivedAt))
	binary.BigEndian.PutUint64(b[24:], math.Float64bits(t.Price))
	binary.BigEndian.PutUint32(b[32:], t.Size)
	copy(b[36:], t.Conditions[:])
	copy(b[40:], t.Symbol[:])
	b[51] = t.Exchange
	b[52] = t.Tape
	return b
}

func UnmarshalTrade(b []byte) *raw.Trade {
	t := &raw.Trade{
		Id:         binary.BigEndian.Uint64(b[0:]),
		Timestamp:  int64(binary.BigEndian.Uint64(b[8:])),
		ReceivedAt: int64(binary.BigEndian.Uint64(b[16:])),
		Price:      math.Float64frombits(binary.BigEndian.Uint64(b[24:])),
		Size:       binary.BigEndian.Uint32(b[32:]),
		Exchange:   b[51],
		Tape:       b[52],
	}
	copy(t.Conditions[:], b[36:40])
	copy(t.Symbol[:], b[40:51])
	return t
}
