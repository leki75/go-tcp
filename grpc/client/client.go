package client

import (
	"context"
	"io"
	"log"

	"github.com/leki75/go-tcp/schema/proto"
	"google.golang.org/grpc"
)

type client struct {
	ch chan<- *proto.Trade
}

func NewClient(ch chan<- *proto.Trade) *client {
	client := &client{
		ch: ch,
	}
	return client
}

func (c *client) Connect(addr string) error {
	log.Println("Connecting to", addr)

	var opts = []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	client := proto.NewStocksClient(conn)
	stream, err := client.Trades(context.Background(), &proto.TradeParams{})
	if err != nil {
		return err
	}

	for {
		trade, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		c.ch <- trade
	}
	return nil
}
