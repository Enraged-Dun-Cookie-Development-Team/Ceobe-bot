package grpc_impl

import (
	"ceobe-bot/Ceobe_Proto/code_gen/pb"
	"ceobe-bot/conf"
	"ceobe-bot/global"
	"context"

	"github.com/tencent-connect/botgo/dto"
)

type LogServer struct {
	pb.UnimplementedLogServer
}

func NewLogServer() *LogServer {
	return new(LogServer)
}

// 从grpc接收日志并推送到频道中
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
	content += "日志等级：" + in.Level.String() + "\n"
	if in.Manual {
		content += "是否人工介入：是\n"
	} else {
		content += "是否人工介入：否\n"
	}
	content += "信息：" + in.Info
	if in.Extra != "" {
		content += "\n---- 以下是额外内容 ----\n" + in.Extra
	}

	// 1000用于改变推送模式，否则0-6点无法主动推送
	toCreate := &dto.MessageToCreate{
		MsgID:   "1000",
		Content: content,
	}
	_, err := global.BOT_PROCESS.Api.PostMessage(ctx, conf.GetConfig().Bot.ChannelNotice, toCreate)
	if err != nil {
		return &pb.LogResponse{Success: false}, nil
	}
	return &pb.LogResponse{Success: true}, nil
}
