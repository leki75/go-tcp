package config

const (
	EncodingBinary = iota + 1
	EncodingText
)

const (
	ProtoTCP = iota + 1
	ProtoGRPC
	ProtoGRPCX
)

var (
	Proto    = ProtoGRPC
	Encoding = EncodingText
)
