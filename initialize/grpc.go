package initialize

import (
	"context"
	"log"
	"net"

	"ceobe-bot/global"
	"ceobe-bot/pb"

	"github.com/tencent-connect/botgo/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	toCreate := &dto.MessageToCreate{
		Content: "主动推送测试",
	}
	global.BOT_PROCESS.Api.PostMessage(ctx, "99368078", toCreate)
	return &pb.HelloResponse{Reply: "success"}, nil
}

func InitGrpc() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("listen err")
	}

	s := grpc.NewServer()

	reflection.Register(s)
	pb.RegisterGreeterServer(s, &greeterServer{})
	log.Println("Serving gRPC on 127.0.0.1" + ":8000")
	err2 := s.Serve(l)
	if err2 != nil {
		log.Fatal("serve err")
	}
}
