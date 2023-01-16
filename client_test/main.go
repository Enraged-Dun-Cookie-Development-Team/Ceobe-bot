package main

import (
	"ceobe-bot/pb"
	"context"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("连接 gPRC 服务失败,", err)
	}

	defer conn.Close()

	// 创建 gRPC 客户端
	grpcClient := pb.NewGreeterClient(conn)

	// 创建请求参数
	request := pb.HelloRequest{
		Name: `test`,
	}

	// 发送请求，调用 MyTest 接口
	response, err := grpcClient.SayHello(context.Background(), &request)
	if err != nil {
		log.Fatal("发送请求失败，原因是:", err)
	}
	log.Println(response)
}
