package json

import (
	"strconv"
	"time"

	"github.com/leki75/go-tcp/schema"
)

func MarshalTrade(t *schema.Trade) ([]byte, error) {
	b := make([]byte, 0, 222)

	// Type
	b = append(b, `{"T":"t"`...) // 8

	// Symbol
	b = append(b, `,"S":"`...) // 14
	i := 0
	for ; i < 11; i++ {
		if t.Symbol[i] == 0 {
			break
		}
	}
	b = append(b, t.Symbol[:i]...) // 25
	b = append(b, '"')             // 26

	// ID
	b = append(b, `,"i":`...)           // 31
	b = strconv.AppendUint(b, t.Id, 10) // 51 - max 20

	// Exchange
	b = append(b, `,"x":"`...)     // 57
	b = append(b, t.Exchange, '"') // 59

	// Price
	b = append(b, `,"p":`...)                        // 64
	b = strconv.AppendFloat(b, t.Price, 'f', -1, 64) // 91 - max 27

	// Size
	b = append(b, `,"s":`...)                     // 96
	b = strconv.AppendUint(b, uint64(t.Size), 10) // 106 - max 10

	// Conditions
	b = append(b, `,"c":[`...)               // 112
	b = append(b, '"', t.Conditions[0], '"') // 115
	i = 1
	for ; i < 4; i++ {
		if t.Conditions[i] != 0 {
			b = append(b, ',', '"', t.Conditions[i], '"') // 127 - 3*4
		}
	}
	b = append(b, ']') // 128

	// Tape
	b = append(b, `,"z":"`...) // 134
	b = append(b, t.Tape, '"') // 136

	// Timestamp
	b = append(b, `,"t":"`...)                                      // 142
	b = time.Unix(0, t.Timestamp).AppendFormat(b, time.RFC3339Nano) // 177 - max 35
	b = append(b, '"', '}')                                         // 179

	// ReceivedAt
	b = append(b, `,"r":"`...)                                       // 185
	b = time.Unix(0, t.ReceivedAt).AppendFormat(b, time.RFC3339Nano) // 220 - max 35
	b = append(b, '"', '}')                                          // 222

	return b, nil
}
