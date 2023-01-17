package initialize

import (
	"log"
	"net"
	"strconv"

	"ceobe-bot/conf"
	"ceobe-bot/grpc_impl"
	"ceobe-bot/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func InitGrpc() {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(conf.GetConfig().Grpc.Port))
	if err != nil {
		log.Fatal("listen err")
	}

	s := grpc.NewServer()

	reflection.Register(s)

	// 注册grpc服务
	pb.RegisterLogServer(s, grpc_impl.NewLogServer())

	log.Println("Serving gRPC on 127.0.0.1:" + strconv.Itoa(conf.GetConfig().Grpc.Port))
	err2 := s.Serve(l)
	if err2 != nil {
		log.Fatal("serve err")
	}
}
