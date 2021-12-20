package main

import (
	"fmt"

	"github.com/leki75/go-tcp/config"
	grpcclient "github.com/leki75/go-tcp/grpc/client"
	"github.com/leki75/go-tcp/schema/proto"
	"github.com/leki75/go-tcp/schema/raw/binary"
	"github.com/leki75/go-tcp/schema/raw/json"
	"github.com/leki75/go-tcp/schema/raw/msgpack"
	tcpclient "github.com/leki75/go-tcp/tcp/client"
	tcpserver "github.com/leki75/go-tcp/tcp/server"
)

func main() {
	var err error

	switch config.Proto {
	case config.TCP:
		in := make(chan []byte, 10000)

		switch config.Encoding {
		case config.Binary:
			go tcpserver.NewServer(in).Listen(":8000")
			err = tcpclient.NewClient(in).Connect(":8080")

		case config.JSON:
			ch := binary.NewReader(in)
			out := make(chan []byte, 128)
			go func() {
				for {
					b := <-ch
					if len(b) == 54 {
						out <- json.MarshalTrade(binary.UnmarshalTrade(b))
					} else {
						fmt.Println(len(b))
					}
				}
			}()
			go tcpserver.NewServer(out).Listen(":8000")
			err = tcpclient.NewClient(in).Connect(":8080")

		case config.MsgPack:
			ch := binary.NewReader(in)
			out := make(chan []byte, 128)
			go func() {
				for {
					b := <-ch
					if len(b) == 54 {
						out <- msgpack.MarshalTrade(binary.UnmarshalTrade(b))
					} else {
						fmt.Println(len(b))
					}
				}
			}()
			go tcpserver.NewServer(out).Listen(":8000")
			err = tcpclient.NewClient(in).Connect(":8080")
		}

	case config.GRPC:
		ch := make(chan *proto.Trade, 10000)
		go tcpserver.NewProtoServer(ch).Listen(":8000")
		err = grpcclient.NewClient(ch).Connect(":9090")
	}

	panic(err)
}
