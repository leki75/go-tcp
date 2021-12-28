package config

type encoding int

const (
	Binary encoding = iota + 1
	MsgPack
	JSON
)

type proto int

const (
	TCP proto = iota + 1
	GRPC
	WebSocket
)

var (
	Encoding = MsgPack
	Proto    = TCP
)

var WebsocketLibrary = "gorilla"
