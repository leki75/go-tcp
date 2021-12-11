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

	buf := make([]byte, 32768)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		c.ch <- buf[:n]
	}
}
