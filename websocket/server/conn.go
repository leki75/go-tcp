package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/leki75/go-tcp/config"
)

// Conn represents a websocket connection between the server and the client
type Conn interface {
	Close() error
	WriteMessage(context.Context, []byte) error
}

var (
	writeWait = 1 * time.Second // Time allowed to write a message to the peer
)

// ConnFactory can be used to create new Conn instances
type ConnFactory struct {
	method func(w http.ResponseWriter, r *http.Request) (Conn, error)
}

// NewConnFactory creates a new ConnFactory based on the configuration
func NewConnFactory() (*ConnFactory, error) {
	f := ConnFactory{}
	switch config.WebsocketLibrary {
	case "gobwas":
		f.method = NewGobwasWebsocketConn
	case "gorilla":
		f.method = NewGorillaWebsocketConn
	case "nhooyr":
		f.method = NewNhooyrWebsocketConn
	default:
		return nil, errors.New("unsupported websocket library: " + config.WebsocketLibrary)
	}
	log.Println("using websocket library: " + config.WebsocketLibrary)
	return &f, nil
}

// NewConn returns a new websocket connection
func (f *ConnFactory) NewConn(w http.ResponseWriter, r *http.Request) (Conn, error) {
	return f.method(w, r)
}
