package msgpack

import (
	"encoding/binary"
	"math"

	"github.com/leki75/go-tcp/schema/proto"
)

func MarshalTrade(t *proto.Trade) []byte {
	var i int
	b := make([]byte, 0, 90)

	// Type - 4
	b = append(b, 0xa1, 'T', 0xa1, 't')

	// Symbol - max 14
	b = append(b, 0xa1, 'S', 0xa0+byte(len(t.Symbol)))
	b = append(b, t.Symbol...)

	// Id - 11
	b = append(b, 0xa1, 'i', 0xcf)
	i = len(b)
	b = append(b, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(b[i:], t.Id)

	// Exchange - 4
	b = append(b, 0xa1, 'x', 0xa1, byte(t.Exchange))

	// Price - 11
	b = append(b, 0xa1, 'p', 0xcb)
	i = len(b)
	b = append(b, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(b[i:], math.Float64bits(t.Price))

	// Size - 7
	b = append(b, 0xa1, 's', 0xce)
	i = len(b)
	b = append(b, 0, 0, 0, 0)
	binary.BigEndian.PutUint32(b[i:], t.Volume)

	// Conditions - max 11
	b = append(b, 0xa1, 'c', 0x90+byte(len(t.Conditions)))
	for i := range t.Conditions {
		b = append(b, 0xa1, t.Conditions[i])
	}

	// Tape - 4
	b = append(b, 0xa1, 'z', 0xa1, byte(t.Tape))

	// Timestamp - 12
	b = append(b, 0xa1, 't', 0xd7, 0xff)
	i = len(b)
	b = append(b, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(b[i:], (uint64(t.Timestamp%1e9)<<34)|uint64(t.Timestamp/1e9))

	// ReceivedAt - 12
	b = append(b, 0xa1, 'r', 0xd7, 0xff)
	i = len(b)
	b = append(b, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(b[i:], (uint64(t.ReceivedAt%1e9)<<34)|uint64(t.ReceivedAt/1e9))

	return b
}
