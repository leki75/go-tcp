package producer

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/leki75/go-tcp/config"
)

// Trades compact binary representation is 45 bytes
type Trade struct {
	Id         uint64   // i - 8
	Timestamp  int64    // t - 16
	Price      float64  // p - 24
	Size       uint32   // s - 28
	Conditions [4]byte  // c - 32
	Symbol     [11]byte // S - 43
	Exchange   byte     // x - 44
	Tape       byte     // z - 45
}

var _ json.Marshaler = (*Trade)(nil)

func (t *Trade) MarshalBinary() ([]byte, error) {
	b := make([]byte, 45)
	binary.BigEndian.PutUint64(b[0:], t.Id)
	binary.BigEndian.PutUint64(b[8:], uint64(t.Timestamp))
	binary.BigEndian.PutUint64(b[16:], math.Float64bits(t.Price))
	binary.BigEndian.PutUint32(b[24:], t.Size)
	for i := 0; i < 4; i++ {
		b[i+28] = t.Conditions[i]
	}
	for i := 0; i < 11; i++ {
		b[i+32] = t.Symbol[i]
	}
	b[43] = t.Exchange
	b[44] = t.Tape
	return b, nil
}

// JSON encoded version
func (t *Trade) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 179)

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

	return b, nil
}

// TradeProducer is blocking should be running in a goroutine
func TradeProducer(ch chan<- []byte) {
	trade := Trade{
		Symbol:     [11]byte{'A', 'A', 'P', 'L'},
		Price:      123.456,
		Size:       100,
		Conditions: [4]byte{'@'},
		Exchange:   'N',
		Tape:       'C',
	}

	start := time.Now()
	processed := 0
	i := uint64(0)

	var b []byte
	for {
		now := time.Now()
		trade.Id = i
		trade.Timestamp = now.UnixNano()

		switch config.Encoding {
		case config.Binary:
			b, _ = trade.MarshalBinary()
		case config.Text:
			b, _ = trade.MarshalJSON()
		}

		if now.Sub(start) >= time.Second {
			log.Println("Processed:", processed)
			start = now
			processed = 0
		}

		ch <- b
		i++
		processed++
	}
}
