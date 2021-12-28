package server

import (
	"log"
	"net/http"
)

var factory *ConnFactory

func init() {
	var err error
	factory, err = NewConnFactory()
	if err != nil {
		panic(err)
	}
}

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

func (s *server) handleConnection(conn Conn, remoteAddr string) {
	s.register <- newClient(s, conn, remoteAddr)
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := factory.NewConn(w, r)
	if err != nil {
		panic(err)
	}

	log.Println("client connect", r.RemoteAddr)
	go s.handleConnection(conn, r.RemoteAddr)
}

func (s *server) Listen(addr string) error {
	log.Println("Listening on", addr)
	http.HandleFunc("/", s.handler)
	return http.ListenAndServe(":8000", nil)
}
