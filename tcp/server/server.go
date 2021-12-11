package server

import (
	"log"
	"net"

	"github.com/leki75/go-tcp/config"
	"github.com/leki75/go-tcp/proto"
	protobinary "github.com/leki75/go-tcp/proto/binary"
	protojson "github.com/leki75/go-tcp/proto/json"
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

func NewProtoServer(ch <-chan *proto.Trade) *server {
	s := &server{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]struct{}),
	}
	go s.runProto(ch)
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

func (s *server) runProto(ch <-chan *proto.Trade) {
	var b []byte
	for {
		select {
		case c := <-s.register:
			s.clients[c] = struct{}{}
		case c := <-s.unregister:
			delete(s.clients, c)
		case trade := <-ch:
			switch config.Encoding {
			case config.EncodingBinary:
				b = protobinary.MarshalTrade(trade)
			case config.EncodingText:
				b = protojson.MarshalTrade(trade)
			}

			for c := range s.clients {
				c.Send(b)
			}
		}
	}
}

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
