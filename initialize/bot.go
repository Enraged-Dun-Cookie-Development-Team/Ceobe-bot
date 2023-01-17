package initialize

import (
	"ceobe-bot/global"
	"ceobe-bot/process"
	"context"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

type BotConfig struct {
	AppId         int64  `json:"app_id" env:"APP_ID,notEmpty" validate:"required"`
	Token         string `json:"token" env:"TOKEN,notEmpty" validate:"required"`
	ChannelNotice string `json:"channel_notice" env:"CHANNEL_NOTICE,notEmpty" validate:"required"`
}

func InitBot(config BotConfig) {
	botToken := token.BotToken(uint64(config.AppId), config.Token)
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 沙箱环境
	// api := botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)
	ctx := context.Background()

	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	global.BOT_PROCESS = &process.Processor{Api: api}

	intent := websocket.RegisterHandlers(
		// at 机器人事件，目前是在这个事件处理中有逻辑，会回消息，其他的回调处理都只把数据打印出来，不做任何处理
		ATMessageEventHandler(),
		// 如果想要捕获到连接成功的事件，可以实现这个回调
		ReadyHandler(),
		// 连接关闭回调
		ErrorNotifyHandler(),
		// 频道事件
		GuildEventHandler(),
		// 成员事件
		MemberEventHandler(),
		// 子频道事件
		ChannelEventHandler(),
		// 私信，目前只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		DirectMessageHandler(),
		// 频道消息，只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		CreateMessageHandler(),
		// 互动事件
		InteractionHandler(),
	)

	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

// ReadyHandler 自定义 ReadyHandler 感知连接成功事件
func ReadyHandler() event.ReadyHandler {
	return func(event *dto.WSPayload, data *dto.WSReadyData) {
		log.Println("ready event receive: ", data)
	}
}

func ErrorNotifyHandler() event.ErrorNotifyHandler {
	return func(err error) {
		log.Println("error notify receive: ", err)
	}
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return global.BOT_PROCESS.ProcessMessage(input, data)
	}
}

// GuildEventHandler 处理频道事件
func GuildEventHandler() event.GuildEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildData) error {
		fmt.Println(data)
		return nil
	}
}

// ChannelEventHandler 处理子频道事件
func ChannelEventHandler() event.ChannelEventHandler {
	return func(event *dto.WSPayload, data *dto.WSChannelData) error {
		fmt.Println(data)
		return nil
	}
}

// MemberEventHandler 处理成员变更事件
func MemberEventHandler() event.GuildMemberEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		fmt.Println(data)
		return nil
	}
}

// DirectMessageHandler 处理私信事件
func DirectMessageHandler() event.DirectMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		fmt.Println(data)
		return nil
	}
}

// CreateMessageHandler 处理消息事件
func CreateMessageHandler() event.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		fmt.Println(data)
		return nil
	}
}

// InteractionHandler 处理内联交互事件
func InteractionHandler() event.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		fmt.Println(data)
		return global.BOT_PROCESS.ProcessInlineSearch(data)
	}
}

func getConfigPath(name string) string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s/%s", path.Dir(filename), name)
	}
	return ""
}
