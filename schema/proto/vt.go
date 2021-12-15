package proto

import (
	"fmt"

	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

type vtCodec struct{}

type VTMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

func (vtCodec) Marshal(v interface{}) ([]byte, error) {
	vt, ok := v.(VTMessage)
	if ok {
		return vt.MarshalVT()
	}

	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return proto.Marshal(vv)
}

func (vtCodec) Unmarshal(data []byte, v interface{}) error {
	vt, ok := v.(VTMessage)
	if ok {
		return vt.UnmarshalVT(data)
	}

	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv)
}

// Name is the name registered for the proto compressor.
func (vtCodec) Name() string {
	return "proto"
}

func init() {
	encoding.RegisterCodec(vtCodec{})
}
