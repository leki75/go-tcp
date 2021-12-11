package binary

import (
	"encoding/binary"
	"math"

	"github.com/leki75/go-tcp/schema"
)

func MarshalTrade(t *schema.Trade) ([]byte, error) {
	b := make([]byte, 53)
	binary.BigEndian.PutUint64(b[0:], t.Id)
	binary.BigEndian.PutUint64(b[8:], uint64(t.Timestamp))
	binary.BigEndian.PutUint64(b[16:], uint64(t.ReceivedAt))
	binary.BigEndian.PutUint64(b[24:], math.Float64bits(t.Price))
	binary.BigEndian.PutUint32(b[32:], t.Size)
	copy(b[36:], t.Conditions[:])
	copy(b[40:], t.Symbol[:])
	b[51] = t.Exchange
	b[52] = t.Tape
	return b, nil
}
