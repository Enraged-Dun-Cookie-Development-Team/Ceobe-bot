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

type LogServer struct {
	pb.UnimplementedLogServer
}

func (l *LogServer) PushLog(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	content := ""
	switch in.Server {
	case pb.LogRequest_RUST:
		content += "服务端：rust端\n"
	case pb.LogRequest_FETCHER:
		content += "服务端：蹲饼器\n"
	case pb.LogRequest_ANALYZER:
		content += "服务端：分析器\n"
	case pb.LogRequest_SCHEDULER:
		content += "服务端：调度器\n"
	}
	content += "日志等级：" + in.Type.String() + "\n"
	if in.Manual {
		content += "是否人工介入：是\n"
	} else {
		content += "是否人工介入：否\n"
	}
	content += "信息：" + in.Info
	if in.Extra != "" {
		content += "\n---- 以下是额外内容 ----\n" + in.Extra
	}

	toCreate := &dto.MessageToCreate{
		Content: content,
	}
	_, err := global.BOT_PROCESS.Api.PostMessage(ctx, "99368078", toCreate)
	if err != nil {
		return &pb.LogResponse{Success: false}, nil
	}
	return &pb.LogResponse{Success: true}, nil
}

func InitGrpc() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("listen err")
	}

	s := grpc.NewServer()

	reflection.Register(s)
	pb.RegisterLogServer(s, &LogServer{})
	log.Println("Serving gRPC on 127.0.0.1" + ":8000")
	err2 := s.Serve(l)
	if err2 != nil {
		log.Fatal("serve err")
	}
}
