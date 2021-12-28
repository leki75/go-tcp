package server

import (
	"context"
	"log"
)

const clientBuffer = 64

type client struct {
	server     *server
	conn       Conn
	remoteAddr string
	data       chan []byte
}

func newClient(s *server, conn Conn, remoteAddr string) *client {
	c := &client{
		server:     s,
		conn:       conn,
		remoteAddr: remoteAddr,
		data:       make(chan []byte, clientBuffer),
	}
	go c.writeLoop()
	return c
}

func (c *client) Send(b []byte) {
	c.data <- b
}

func (c *client) writeLoop() {
	defer c.close()
	for message := range c.data {
		err := c.conn.WriteMessage(context.Background(), message)
		if err != nil {
			break
		}
	}
}

func (c *client) close() {
	go func() {
		c.server.unregister <- c
		close(c.data)
	}()
	for range c.data {
	}
	_ = c.conn.Close()
	log.Println("client close", c.remoteAddr)
}
