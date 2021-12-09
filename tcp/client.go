package tcp

import (
	"encoding/binary"
	"net"

	"github.com/leki75/go-tcp/config"
)

const clientBuffer = 128

type client struct {
	server    *server
	conn      net.Conn
	data      chan []byte
	writeable bool
}

func newClient(s *server, conn net.Conn) *client {
	c := &client{
		server: s,
		conn:   conn,
		data:   make(chan []byte, clientBuffer),
	}
	return c
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

func (c *client) write() {
	var buf net.Buffers
	select {
	case message := <-c.data:
		switch config.Encoding {
		case config.Binary:
			buf = binaryMessage(message, c.data)
		case config.Text:
			buf = jsonMessage(message, c.data)
		}

		_, err := buf.WriteTo(c.conn)
		if err != nil {
			c.close()
		}
		c.writeable = false

	default:
		c.writeable = true
	}

}

func (c *client) close() {
	go func() {
		_ = c.server.removeConnection(c.conn)
		close(c.data)
	}()
	for range c.data {
	}
	_ = c.conn.Close()
}
