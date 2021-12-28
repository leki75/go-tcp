package server

import (
	"context"
	"net/http"

	"nhooyr.io/websocket"
)

// NhooyrWebsocketConn is a websocket connection using the nhooyr websocket library:
// https://pkg.go.dev/nhooyr.io/websocket
type NhooyrWebsocketConn struct {
	conn *websocket.Conn
}

// NewNhooyrWebsocketConn creates a new nhooyr websocket connection
func NewNhooyrWebsocketConn(w http.ResponseWriter, r *http.Request) (Conn, error) {
	opts := websocket.AcceptOptions{
		// NOTE: disable validation of HTTP Origin header
		InsecureSkipVerify: true,
	}
	opts.CompressionMode = websocket.CompressionDisabled
	c, err := websocket.Accept(w, r, &opts)
	if err != nil {
		return nil, err
	}
	return &NhooyrWebsocketConn{
		conn: c,
	}, nil
}

// Close closes the websocket connection
func (c *NhooyrWebsocketConn) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "")
}

// WriteMessage writes a single message
func (c *NhooyrWebsocketConn) WriteMessage(ctx context.Context, data []byte) error {
	writeCtx, cancel := context.WithTimeout(ctx, writeWait)
	defer cancel()

	return c.conn.Write(writeCtx, websocket.MessageBinary, data)
}
