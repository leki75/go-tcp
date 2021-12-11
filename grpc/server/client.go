package server

import (
	"github.com/leki75/go-tcp/proto"
)

type stocksClient struct {
	trades chan *proto.Trade
}

func newStocksClient() *stocksClient {
	c := &stocksClient{
		trades: make(chan *proto.Trade),
	}
	return c
}

func (c *stocksClient) Send(trade *proto.Trade) {
	c.trades <- trade
}
