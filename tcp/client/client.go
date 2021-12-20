package client

import (
	"log"
	"net"
)

type client struct {
	ch chan<- []byte
}

func NewClient(ch chan<- []byte) *client {
	client := &client{
		ch: ch,
	}
	return client
}

func (c *client) Connect(addr string) error {
	log.Println("Connecting to", addr)
	defer close(c.ch)

	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return err
	}

	buf := make([]byte, 65536)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		b := make([]byte, n)
		copy(b, buf[:n])
		c.ch <- b
	}
}
