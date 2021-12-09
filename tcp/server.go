package tcp

import (
	"log"
	"net"
)

type server struct {
	register   chan *client
	unregister chan *client
	clients    map[*client]struct{}
}

func NewServer(ch <-chan []byte) *server {
	server := &server{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]struct{}),
	}
	go server.run(ch)
	return server
}

func (s *server) run(ch <-chan []byte) {
	for {
		select {
		case c := <-s.register:
			s.clients[c] = struct{}{}
		case c := <-s.unregister:
			delete(s.clients, c)
		case b := <-ch:
			for c := range s.clients {
				c.Send(b)
			}
		}
	}
}

// Listen is blocking
func (s *server) Listen(addr string) error {
	log.Println("Listening on", addr)
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
	log.Println("client connect", conn.RemoteAddr())
	s.register <- newClient(s, conn)
}
