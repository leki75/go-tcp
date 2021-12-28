package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// GorillaWebsocketConn is a websocket connection using the gorilla websocket library:
// https://pkg.go.dev/github.com/gorilla/websocket
type GorillaWebsocketConn struct {
	conn *websocket.Conn
}

// NewGorillaWebsocketConn creates a new gorilla websocket connection
func NewGorillaWebsocketConn(w http.ResponseWriter, r *http.Request) (Conn, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &GorillaWebsocketConn{
		conn: c,
	}, nil
}

// Close closes the websocket connection
func (c *GorillaWebsocketConn) Close() error {
	return c.conn.Close()
}

// WriteMessage writes a single message
func (c *GorillaWebsocketConn) WriteMessage(ctx context.Context, data []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
