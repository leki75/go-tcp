package server

import (
	"log"
	"net"

	"github.com/leki75/go-tcp/schema/proto"
	"google.golang.org/grpc"
)

type stocksServer struct {
	proto.UnimplementedStocksServer

	register   chan *stocksClient
	unregister chan *stocksClient
	clients    map[*stocksClient]struct{}
}

func NewServer(ch <-chan *proto.Trade) *stocksServer {
	s := &stocksServer{
		register:   make(chan *stocksClient),
		unregister: make(chan *stocksClient),
		clients:    make(map[*stocksClient]struct{}),
	}
	go s.run(ch)
	return s
}

func (s *stocksServer) run(ch <-chan *proto.Trade) {
	for {
		select {
		case c := <-s.register:
			s.clients[c] = struct{}{}
		case c := <-s.unregister:
			delete(s.clients, c)
		case trade := <-ch:
			for c := range s.clients {
				c.Send(trade)
			}
		}
	}
}

func (s *stocksServer) Listen(addr string) error {
	log.Println("Listening on", addr)
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	var opts = []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterStocksServer(grpcServer, s)

	return grpcServer.Serve(listener)
}

func (s *stocksServer) Trades(_ *proto.TradeParams, stream proto.Stocks_TradesServer) error {
	ctx := stream.Context()
	c := newStocksClient()

	s.register <- c
	defer func() {
		s.unregister <- c
	}()

	for {
		select {
		case trade := <-c.trades:
			if err := stream.Send(trade); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
