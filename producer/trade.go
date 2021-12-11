package producer

import (
	"log"
	"time"

	"github.com/leki75/go-tcp/config"
	"github.com/leki75/go-tcp/proto"
	protobinary "github.com/leki75/go-tcp/proto/binary"
	"github.com/leki75/go-tcp/schema"
	schemabinary "github.com/leki75/go-tcp/schema/binary"
	schemajson "github.com/leki75/go-tcp/schema/json"
)

// TCPTradeProducer is blocking should be running in a goroutine
func TCPTradeProducer(ch chan<- []byte) {
	var b []byte
	tradeProducer(func(t *schema.Trade) {
		switch config.Encoding {
		case config.EncodingBinary:
			b, _ = schemabinary.MarshalTrade(t)
		case config.EncodingText:
			b, _ = schemajson.MarshalTrade(t)
		}
		ch <- b
	})
}

// GRPCTradeProducer is blocking should be running in a goroutine
func GRPCTradeProducer(ch chan<- *proto.Trade) {
	tradeProducer(func(t *schema.Trade) {
		ch <- &proto.Trade{
			Id:         t.Id,
			Timestamp:  uint64(t.Timestamp),
			ReceivedAt: t.ReceivedAt,
			Price:      t.Price,
			Volume:     t.Size,
			Conditions: protobinary.GetConditions(t.Conditions[:]),
			Symbol:     protobinary.GetSymbol(t.Symbol[:]),
			Exchange:   int32(t.Exchange),
			Tape:       int32(t.Tape),
		}
	})
}

func tradeProducer(fn func(*schema.Trade)) {
	trade := schema.Trade{
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

		fn(&trade)

		if now.Add(-time.Second).After(start) {
			log.Println("processed", processed)
			start = now
			processed = 0
		}

		i++
		processed++
	}
}
