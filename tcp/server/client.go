package server

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/leki75/go-tcp/config"
)

const clientBuffer = 64

type client struct {
	server *server
	conn   net.Conn
	data   chan []byte
}

func newClient(s *server, conn net.Conn) *client {
	c := &client{
		server: s,
		conn:   conn,
		data:   make(chan []byte, clientBuffer),
	}
	go c.writeLoop()
	return c
}

func (c *client) Send(b []byte) {
	// TODO: wrap in a select to not block the main routine
	c.data <- b
}

var (
	openingBracket = []byte{'[', '\n'}
	separatorComma = []byte{',', '\n'}
	closingBracket = []byte{'\n', ']', '\n'}
)

func binaryMessage(message []byte, ch <-chan []byte) net.Buffers {
	size := len(ch)
	buf := make(net.Buffers, size+2)
	buf[0] = []byte{'t', 0, 0} // header: type(1), size(2)
	binary.BigEndian.PutUint16(buf[0][1:], uint16(size+1))
	buf[1] = message
	size += 2
	for i := 2; i < size; i++ {
		buf[i] = <-ch
	}
	return buf
}

func jsonMessage(message []byte, ch <-chan []byte) net.Buffers {
	size := len(ch)
	buf := make(net.Buffers, size*2+3)
	buf[0] = openingBracket
	buf[1] = message
	buf[size*2+2] = closingBracket
	for i, j := 0, 2; i < size; i++ {
		buf[j] = separatorComma
		j++
		buf[j] = <-ch
		j++
	}
	return buf
}

func (c *client) writeLoop() {
	defer c.close()

	var buf net.Buffers
	for message := range c.data {
		switch config.Encoding {
		case config.EncodingBinary:
			buf = binaryMessage(message, c.data)
		case config.EncodingText:
			buf = jsonMessage(message, c.data)
		}

		_, err := buf.WriteTo(c.conn)
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
	log.Println("client close", c.conn.RemoteAddr())
}
