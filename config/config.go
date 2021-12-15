package config

const (
	EncodingBinary = iota + 1
	EncodingMsgPack
	EncodingJSON
)

const (
	ProtoTCP = iota + 1
	ProtoGRPC
	ProtoWS
)

var (
	Encoding = EncodingJSON
	Proto    = ProtoTCP
)
