package binary

import (
	"encoding/binary"
	"fmt"
)

const (
	stateHeader = iota
	stateBody
)

type reader struct {
	state byte
	count uint16
	buf   []byte
	in    <-chan []byte
	out   chan []byte
}

func NewReader(in <-chan []byte) <-chan []byte {
	r := &reader{
		in:  in,
		out: make(chan []byte, 128),
	}
	go r.process()
	return r.out
}

func (r *reader) process() {
	var b []byte
	for {
		b = <-r.in
	header:
		for {
			switch r.state {
			case stateHeader:
				size := 2 - len(r.buf)
				if len(b) < size {
					r.state = stateHeader
					r.buf = append(r.buf, b...)
					break header
				}

				r.buf = append(r.buf, b[:size]...)
				r.count = binary.BigEndian.Uint16(r.buf)
				b = b[size:]

			case stateBody:
				if len(r.buf) == 0 {
					break
				}

				size := int(r.buf[0]) + 1 - len(r.buf)
				if len(b) < int(size) {
					r.state = stateBody
					r.buf = append(r.buf, b...)
					break header
				}

				r.buf = r.buf[1:]
				r.buf = append(r.buf, b[:size]...)
				b = b[size:]
				r.out <- r.buf
				r.count--
				if r.count == 0 {
					r.state = stateHeader
					r.buf = []byte{}
					continue
				}
			}

			for len(b) > 0 && r.count > 0 {
				if len(b) == 1 {
					r.state = stateBody
					r.buf = b
					break header
				}

				size := b[0]
				if len(b) < int(size)+1 {
					r.state = stateBody
					r.buf = b
					break header
				}
				b = b[1:]

				if b[0] != 't' {
					fmt.Println(len(b), r, b[0])
					panic(size)
				}
				r.out <- b[:size]
				b = b[size:]
				r.count--
			}

			if len(b) == 0 && r.count == 0 {
				r.state = stateHeader
				r.buf = []byte{}
				break header
			} else if len(b) == 0 {
				r.state = stateBody
				r.buf = []byte{}
				break header
			} else if r.count == 0 {
				r.state = stateHeader
				r.buf = []byte{}
			}
		}
	}
}
