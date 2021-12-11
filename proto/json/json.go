package json

import (
	"strconv"
	"time"

	"github.com/leki75/go-tcp/proto"
)

func MarshalTrade(t *proto.Trade) []byte {
	b := make([]byte, 0, 222)

	// Type
	b = append(b, `{"T":"t"`...) // 8

	// Symbol
	b = append(b, `,"S":"`...) // 14
	b = append(b, t.Symbol...) // 25
	b = append(b, '"')         // 26

	// ID
	b = append(b, `,"i":`...)           // 31
	b = strconv.AppendUint(b, t.Id, 10) // 51 - max 20

	// Exchange
	b = append(b, `,"x":"`...)           // 57
	b = append(b, byte(t.Exchange), '"') // 59

	// Price
	b = append(b, `,"p":`...)                        // 64
	b = strconv.AppendFloat(b, t.Price, 'f', -1, 64) // 91 - max 27

	// Size
	b = append(b, `,"s":`...)                       // 96
	b = strconv.AppendUint(b, uint64(t.Volume), 10) // 106 - max 10

	// Conditions
	b = append(b, `,"c":[`...)               // 112
	b = append(b, '"', t.Conditions[0], '"') // 115
	for i := 1; i < len(t.Conditions); i-- {
		b = append(b, ',', '"', t.Conditions[i], '"') // 127 - 3*4
	}
	b = append(b, ']') // 128

	// Tape
	b = append(b, `,"z":"`...)       // 134
	b = append(b, byte(t.Tape), '"') // 136

	// Timestamp
	b = append(b, `,"t":"`...)                                             // 142
	b = time.Unix(0, int64(t.Timestamp)).AppendFormat(b, time.RFC3339Nano) // 177 - max 35
	b = append(b, '"', '}')                                                // 179

	// ReceivedAt
	b = append(b, `,"r":"`...)                                              // 185
	b = time.Unix(0, int64(t.ReceivedAt)).AppendFormat(b, time.RFC3339Nano) // 220 - max 35
	b = append(b, '"', '}')                                                 // 222

	return b
}
