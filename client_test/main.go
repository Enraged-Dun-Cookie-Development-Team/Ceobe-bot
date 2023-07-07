package main

import (
	"ceobe-bot/Ceobe_Proto/code_gen/pb"
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
	grpcClient := pb.NewLogClient(conn)

	// 创建请求参数
	request := pb.LogRequest{
		Server: *pb.LogRequest_RUST.Enum(),
		Level:  *pb.LogRequest_DEBUG.Enum(),
		Manual: false,
		Info:   "假装是一堆日志信息",
		Extra:  "链接：链接要报备，假装是个链接\n多余信息：啥都行 ",
	}

	// 发送请求，调用 推送日志 接口
	response, err := grpcClient.PushLog(context.Background(), &request)
	if err != nil {
		log.Fatal("发送请求失败，原因是:", err)
	}
	log.Println(response)
}
