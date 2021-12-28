package server

import (
	"log"
	"net"

	"github.com/gobwas/ws"
)

type server struct {
	register   chan *client
	unregister chan *client
	clients    map[*client]struct{}
}

func NewServer(ch <-chan []byte) *server {
	s := &server{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]struct{}),
	}
	go s.run(ch)
	return s
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

func (s *server) Listen(addr string) error {
	log.Println("Listening on gobwas", addr)
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		_, err = ws.Upgrade(conn)
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
