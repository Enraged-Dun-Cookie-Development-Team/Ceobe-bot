package initialize

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type T struct {

}

func InitGrpc() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("listen err")
	}

	s := grpc.NewServer()

	reflection.Register(s)

	err2 := s.Serve(l)
	if err2 != nil {
		log.Fatal("serve err")
	}
}