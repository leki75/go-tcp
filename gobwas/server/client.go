package server

import (
	"log"
	"net"

	"github.com/gobwas/ws"
)

const clientBuffer = 64

type client struct {
	server     *server
	conn       net.Conn
	remoteAddr string
	data       chan []byte
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
	c.data <- b
}

func binaryMessage(message []byte, ch <-chan []byte) (net.Buffers, int) {
	numMessages := len(ch) + 1
	buf := make(net.Buffers, numMessages)
	length := len(message)

	buf[0] = message

	var b []byte
	for i := 1; i < numMessages; i++ {
		b = <-ch
		buf[i] = b
		length += len(b)
	}
	return buf, length
}

func (c *client) writeLoop() {
	defer c.close()

	header := ws.Header{
		Fin:    true,
		OpCode: ws.OpText,
	}
	// {
	// 	Fin    bool
	// 	Rsv    byte
	// 	OpCode OpCode
	// 	Masked bool
	// 	Mask   [4]byte
	// 	Length int64
	// }

	var buf net.Buffers
	var n int
	for message := range c.data {
		buf, n = binaryMessage(message, c.data)
		header.Length = int64(n)
		//log.Println(header)
		if err := ws.WriteHeader(c.conn, header); err != nil {
			break
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
	log.Println("client close", c.remoteAddr)
}
