package producer

import (
	"log"
	"time"

	"github.com/leki75/go-tcp/schema/proto"
	"github.com/leki75/go-tcp/schema/raw"
	rawbinary "github.com/leki75/go-tcp/schema/raw/binary"
)

// ByteSliceTradeProducer is blocking should be running in a goroutine
func ByteSliceTradeProducer(ch chan<- []byte) {
	trade := raw.Trade{
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

	for {
		now := time.Now()
		trade.Id = i
		trade.ReceivedAt = now.UnixNano()
		trade.Timestamp = trade.ReceivedAt - int64(time.Millisecond)

		ch <- rawbinary.MarshalTrade(&trade)

		if now.Add(-time.Second).After(start) {
			log.Println("processed", processed)
			start = now
			processed = 0
		}

		i++
		processed++
	}
}

// ProtoTradeProducer is blocking should be running in a goroutine
func ProtobufTradeProducer(ch chan<- *proto.Trade) {
	trade := proto.Trade{
		Symbol:     "AAPL",
		Price:      123.456,
		Volume:     100,
		Conditions: []byte{'@'},
		Exchange:   'N',
		Tape:       'C',
	}

	start := time.Now()
	processed := 0
	i := uint64(0)

	for {
		now := time.Now()
		trade.Id = i
		trade.ReceivedAt = now.UnixNano()
		trade.Timestamp = uint64(trade.ReceivedAt) - uint64(time.Millisecond)

		ch <- &trade

		if now.Add(-time.Second).After(start) {
			log.Println("processed", processed)
			start = now
			processed = 0
		}

		i++
		processed++
	}
}
