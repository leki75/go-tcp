package server

import (
	"log"
	"net"
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
	numMessages := len(ch) + 1
	buf := make(net.Buffers, numMessages*2+1)

	buf[0] = []byte{byte(numMessages >> 8), byte(numMessages)} // message count header
	// buf[1] = []byte{byte(len(message) >> 8), byte(len(message))}
	buf[1] = []byte{byte(len(message))}
	buf[2] = message
	var b []byte
	count := 0
	for i := 3; i < numMessages*2+1; i += 2 {
		b = <-ch
		//buf[i] = []byte{byte(len(b) >> 8), byte(len(b))}
		buf[i] = []byte{byte(len(b))}
		buf[i+1] = b
		count++
	}
	return buf
}

// func jsonMessage(message []byte, ch <-chan []byte) net.Buffers {
// 	numMessages := len(ch) + 1
// 	buf := make(net.Buffers, numMessages*2+1)
// 	buf[0] = openingBracket
// 	buf[1] = message
// 	buf[numMessages*2] = closingBracket
// 	for i := 2; i < numMessages*2+1; i += 2 {
// 		buf[i] = separatorComma
// 		buf[i+1] = <-ch
// 	}
// 	return buf
// }

func (c *client) writeLoop() {
	defer c.close()

	var buf net.Buffers
	for message := range c.data {
		// switch config.Encoding {
		// case config.EncodingBinary:
		// 	buf = binaryMessage(message, c.data)
		// case config.EncodingJSON:
		// 	buf = jsonMessage(message, c.data)
		// }
		buf = binaryMessage(message, c.data)
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
