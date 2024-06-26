package grpc_impl

import (
	"ceobe-bot/Ceobe_Proto/code_gen/pb"
	"ceobe-bot/conf"
	"ceobe-bot/global"
	"context"
	"time"

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
	now := time.Now()
	// 格式化为精确到秒的时间字符串
	timeStr := now.Format("2006-01-02 15:04:05")
	content += "当前时间：" + timeStr + "\n"
	content += "信息：" + in.Info
	if in.Extra != "" {
		content += "\n---- 以下是额外内容 ----\n" + in.Extra
	}

	// 1000用于改变推送模式，否则0-6点无法主动推送
	toCreate := &dto.MessageToCreate{
		MsgID:   "1000",
		Content: content,
	}
	var err error
	if (in.Level == pb.LogRequest_ERROR || in.Level == pb.LogRequest_WARN) || conf.GetConfig().Bot.ChannelInfo == "" {
		_, err = global.BOT_PROCESS.Api.PostMessage(ctx, conf.GetConfig().Bot.ChannelNotice, toCreate)
	} else {
		_, err = global.BOT_PROCESS.Api.PostMessage(ctx, conf.GetConfig().Bot.ChannelInfo, toCreate)
	}
	if err != nil {
		return &pb.LogResponse{Success: false}, nil
	}
	return &pb.LogResponse{Success: true}, nil
}
