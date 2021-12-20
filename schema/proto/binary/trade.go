package binary

import (
	"encoding/binary"
	"math"

	"github.com/leki75/go-tcp/schema/proto"
)

func MarshalTrade(t *proto.Trade) []byte {
	b := make([]byte, 54)
	b[0] = 't'
	copy(b[1:12], t.Symbol)
	copy(b[12:16], t.Conditions)
	binary.BigEndian.PutUint64(b[16:], t.Id)
	binary.BigEndian.PutUint64(b[24:], uint64(t.Timestamp))
	binary.BigEndian.PutUint64(b[32:], uint64(t.ReceivedAt))
	binary.BigEndian.PutUint64(b[40:], math.Float64bits(t.Price))
	binary.BigEndian.PutUint32(b[48:], t.Volume)
	b[52] = byte(t.Exchange)
	b[53] = byte(t.Tape)
	return b
}

func UnmarshalTrade(b []byte) *proto.Trade {
	return &proto.Trade{
		Id:         binary.BigEndian.Uint64(b[16:]),
		Timestamp:  binary.BigEndian.Uint64(b[24:]),
		ReceivedAt: int64(binary.BigEndian.Uint64(b[32:])),
		Price:      math.Float64frombits(binary.BigEndian.Uint64(b[40:])),
		Volume:     binary.BigEndian.Uint32(b[48:]),
		Symbol:     GetSymbol(b[1:12]),
		Conditions: GetConditions(b[12:16]),
		Exchange:   int32(b[52]),
		Tape:       int32(b[53]),
	}
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
