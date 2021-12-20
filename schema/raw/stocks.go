package raw

import (
	"fmt"
	"time"
)

type Trade struct {
	Id         uint64   // i - 8
	Timestamp  int64    // t - 16
	ReceivedAt int64    // t - 24
	Price      float64  // p - 32
	Size       uint32   // s - 36
	Conditions [4]byte  // c - 40
	Symbol     [11]byte // S - 51
	Exchange   byte     // x - 52
	Tape       byte     // z - 53
}

func (t *Trade) String() string {
	return fmt.Sprintf(
		"Id: %d, Timestamp: %s, Symbol: %s, Conditions: %v, Price: %f, Size: %d",
		t.Id, time.Unix(0, t.Timestamp), t.Symbol, t.Conditions, t.Price, t.Size,
	)
}
