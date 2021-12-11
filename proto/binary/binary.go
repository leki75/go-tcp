package binary

import (
	"encoding/binary"
	"math"

	"github.com/leki75/go-tcp/proto"
)

func MarshalTrade(t *proto.Trade) []byte {
	b := make([]byte, 53)
	binary.BigEndian.PutUint64(b[0:], t.Id)
	binary.BigEndian.PutUint64(b[8:], uint64(t.Timestamp))
	binary.BigEndian.PutUint64(b[16:], uint64(t.ReceivedAt))
	binary.BigEndian.PutUint64(b[24:], math.Float64bits(t.Price))
	binary.BigEndian.PutUint32(b[32:], t.Volume)
	copy(b[36:40], t.Conditions)
	copy(b[40:51], t.Symbol)
	b[51] = byte(t.Exchange)
	b[52] = byte(t.Tape)
	return b
}

func UnmarshalTrade(b []byte) *proto.Trade {
	t := &proto.Trade{
		Id:         binary.BigEndian.Uint64(b),
		Timestamp:  binary.BigEndian.Uint64(b[8:]),
		ReceivedAt: int64(binary.BigEndian.Uint64(b[16:])),
		Price:      math.Float64frombits(binary.BigEndian.Uint64(b[24:])),
		Volume:     binary.BigEndian.Uint32(b[32:]),
		Conditions: GetConditions(b[36:40]),
		Symbol:     GetSymbol(b[40:51]),
		Exchange:   int32(b[51]),
		Tape:       int32(b[52]),
	}
	return t
}

func GetConditions(cond []byte) []byte {
	i := 1
	for ; i < 4; i++ {
		if cond[i] == 0 {
			break
		}
	}
	return cond[:i]
}

func GetSymbol(symbol []byte) string {
	for i, c := range symbol {
		if c == 0 {
			return string(symbol[:i])
		}
	}
	return string(symbol[:])
}
