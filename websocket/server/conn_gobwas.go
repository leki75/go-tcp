package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type GobwasWebsocketConn struct {
	conn net.Conn
}

func NewGobwasWebsocketConn(w http.ResponseWriter, r *http.Request) (Conn, error) {
	c, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return nil, err
	}
	return &GobwasWebsocketConn{
		conn: c,
	}, nil
}

// Close closes the websocket connection
func (c *GobwasWebsocketConn) Close() error {
	return c.conn.Close()
}

// WriteMessage writes a single message
func (c *GobwasWebsocketConn) WriteMessage(ctx context.Context, data []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return wsutil.WriteServerMessage(c.conn, ws.OpBinary, data)
}
