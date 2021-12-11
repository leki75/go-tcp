// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StocksClient is the client API for Stocks service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StocksClient interface {
	Trades(ctx context.Context, in *TradeParams, opts ...grpc.CallOption) (Stocks_TradesClient, error)
}

type stocksClient struct {
	cc grpc.ClientConnInterface
}

func NewStocksClient(cc grpc.ClientConnInterface) StocksClient {
	return &stocksClient{cc}
}

func (c *stocksClient) Trades(ctx context.Context, in *TradeParams, opts ...grpc.CallOption) (Stocks_TradesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Stocks_ServiceDesc.Streams[0], "/gotcp.Stocks/Trades", opts...)
	if err != nil {
		return nil, err
	}
	x := &stocksTradesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Stocks_TradesClient interface {
	Recv() (*Trade, error)
	grpc.ClientStream
}

type stocksTradesClient struct {
	grpc.ClientStream
}

func (x *stocksTradesClient) Recv() (*Trade, error) {
	m := new(Trade)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StocksServer is the server API for Stocks service.
// All implementations must embed UnimplementedStocksServer
// for forward compatibility
type StocksServer interface {
	Trades(*TradeParams, Stocks_TradesServer) error
	mustEmbedUnimplementedStocksServer()
}

// UnimplementedStocksServer must be embedded to have forward compatible implementations.
type UnimplementedStocksServer struct {
}

func (UnimplementedStocksServer) Trades(*TradeParams, Stocks_TradesServer) error {
	return status.Errorf(codes.Unimplemented, "method Trades not implemented")
}
func (UnimplementedStocksServer) mustEmbedUnimplementedStocksServer() {}

// UnsafeStocksServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StocksServer will
// result in compilation errors.
type UnsafeStocksServer interface {
	mustEmbedUnimplementedStocksServer()
}

func RegisterStocksServer(s grpc.ServiceRegistrar, srv StocksServer) {
	s.RegisterService(&Stocks_ServiceDesc, srv)
}

func _Stocks_Trades_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TradeParams)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StocksServer).Trades(m, &stocksTradesServer{stream})
}

type Stocks_TradesServer interface {
	Send(*Trade) error
	grpc.ServerStream
}

type stocksTradesServer struct {
	grpc.ServerStream
}

func (x *stocksTradesServer) Send(m *Trade) error {
	return x.ServerStream.SendMsg(m)
}

// Stocks_ServiceDesc is the grpc.ServiceDesc for Stocks service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stocks_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gotcp.Stocks",
	HandlerType: (*StocksServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Trades",
			Handler:       _Stocks_Trades_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/stocks.proto",
}
