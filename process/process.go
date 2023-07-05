package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

// Processor is a struct to process message
type Processor struct {
	Api openapi.OpenAPI
}

// structure for json format
type MaintainerInfo struct {
	RUST      []string
	FETCHER   []string
	ANALYZER  []string
	SCHEDULER []string
}

// slices of four parts' maintainers
var RustMaintainers, FetcherMaintainers, AnalyzerMaintainers, SchedulerMaintainers []string
var Maintaininfo = []MaintainerInfo{{RustMaintainers, FetcherMaintainers, AnalyzerMaintainers, SchedulerMaintainers}}

// ProcessMessage is a function to process message
func (p *Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	toCreate := &dto.MessageToCreate{
		MsgID:   data.ID,
		Content: "<@" + data.Author.ID + ">默认回复" + message.Emoji(307),
	}

	// 进入到私信逻辑
	if cmd.Cmd == "dm" {
		p.dmHandler(data)
		return nil
	}

	switch cmd.Cmd {
	case "hi":
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "time":
		toCreate.Content = genReplyContent(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "ark":
		toCreate.Ark = genReplyArk(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "公告":
		p.setAnnounces(ctx, data)
	case "pin":
		if data.MessageReference != nil {
			p.setPins(ctx, data.ChannelID, data.MessageReference.MessageID)
		}
	case "emoji":
		if data.MessageReference != nil {
			p.setEmoji(ctx, data.ChannelID, data.MessageReference.MessageID)
		}
	case "添加负责人":
		switch addMaintainer(data) {
		case 1:
			toCreate.Content = "添加成功"
		case 0:
			toCreate.Content = "添加失败，已存在"
		case -1:
			toCreate.Content = "error"
		default:
		}
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "查询负责人":
		toCreate.Content = "该端负责人为：" + searchMaintainer(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "删除负责人":
		switch deleteMaintainer(data) {
		case 1:
			toCreate.Content = "删除成功"
		case 0:
			toCreate.Content = "不存在指定负责人"
		case -1:
			toCreate.Content = "error"
		default:
		}
		p.sendReply(ctx, data.ChannelID, toCreate)

	default:
	}

	return nil
}

// ProcessInlineSearch is a function to process inline search
func (p *Processor) ProcessInlineSearch(interaction *dto.WSInteractionData) error {
	if interaction.Data.Type != dto.InteractionDataTypeChatSearch {
		return fmt.Errorf("interaction data type not chat search")
	}
	search := &dto.SearchInputResolved{}
	if err := json.Unmarshal(interaction.Data.Resolved, search); err != nil {
		log.Println(err)
		return err
	}
	if search.Keyword != "test" {
		return fmt.Errorf("resolved search key not allowed")
	}
	searchRsp := &dto.SearchRsp{
		Layouts: []dto.SearchLayout{
			{
				LayoutType: 0,
				ActionType: 0,
				Title:      "内联搜索",
				Records: []dto.SearchRecord{
					{
						Cover: "https://pub.idqqimg.com/pc/misc/files/20211208/311cfc87ce394c62b7c9f0508658cf25.png",
						Title: "内联搜索标题",
						Tips:  "内联搜索 tips",
						URL:   "https://www.qq.com",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(searchRsp)
	if err := p.Api.PutInteraction(context.Background(), interaction.ID, string(body)); err != nil {
		log.Println("api call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}

func (p *Processor) dmHandler(data *dto.WSATMessageData) {
	dm, err := p.Api.CreateDirectMessage(
		context.Background(), &dto.DirectMessageToCreate{
			SourceGuildID: data.GuildID,
			RecipientID:   data.Author.ID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	toCreate := &dto.MessageToCreate{
		Content: "默认私信回复",
	}
	_, err = p.Api.PostDirectMessage(
		context.Background(), dm, toCreate,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func genReplyContent(data *dto.WSATMessageData) string {
	var tpl = `你好：%s
在子频道 %s 收到消息。
收到的消息发送时时间为：%s
当前本地时间为：%s
消息来自：%s
`

	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		tpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
	)
}

func genReplyArk(data *dto.WSATMessageData) *dto.Ark {
	return &dto.Ark{
		TemplateID: 23,
		KV: []*dto.ArkKV{
			{
				Key:   "#DESC#",
				Value: "这是 ark 的描述信息",
			},
			{
				Key:   "#PROMPT#",
				Value: "这是 ark 的摘要信息",
			},
			{
				Key: "#LIST#",
				Obj: []*dto.ArkObj{
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这里展示的是 23 号模板",
							},
						},
					},
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这是 ark 的列表项名称",
							},
							{
								Key:   "link",
								Value: "https://www.qq.com",
							},
						},
					},
				},
			},
		},
	}
}

// 添加负责人语句format：@bot 添加负责人 端名 @负责人
// eg. @小刻-测试中 添加负责人 RUST @薄生
func addMaintainer(data *dto.WSATMessageData) int {
	var maintainers []*(dto.User) = data.Mentions
	var maintainpart *[]string
	pattern := regexp.MustCompile(`添加负责人\s(\w+)\b`)
	matches := pattern.FindStringSubmatch(data.Content)
	part := ""
	if len(matches) > 1 {
		part = matches[1]
	}
	switch part {
	case "RUST":
		maintainpart = &Maintaininfo[0].RUST
	case "FETCHER":
		maintainpart = &Maintaininfo[0].FETCHER
	case "ANALYZER":
		maintainpart = &Maintaininfo[0].ANALYZER
	case "SCHEDULER":
		maintainpart = &Maintaininfo[0].SCHEDULER
	default:

	}
	filePath := "./conf/maintainers.json"
	file, err := os.Open(filePath)
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return -1
	}
	defer file.Close()
	for _, element := range maintainers[1:] {
		for _, value := range *maintainpart {
			if value == element.ID {
				return 0
			}
		}
		(*maintainpart) = append((*maintainpart), element.ID)
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(Maintaininfo)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return -1
	}
	return 1
}

// 删除负责人语句format：@bot 删除负责人 端名 @负责人
// eg. @小刻-测试中 删除负责人 RUST @薄生
func deleteMaintainer(data *dto.WSATMessageData) int {
	var maintainers []*(dto.User) = data.Mentions
	var maintainpart *[]string
	pattern := regexp.MustCompile(`删除负责人\s(\w+)\b`)
	matches := pattern.FindStringSubmatch(data.Content)
	part := ""
	if len(matches) > 1 {
		part = matches[1]
	}
	switch part {
	case "RUST":
		maintainpart = &Maintaininfo[0].RUST
	case "FETCHER":
		maintainpart = &Maintaininfo[0].FETCHER
	case "ANALYZER":
		maintainpart = &Maintaininfo[0].ANALYZER
	case "SCHEDULER":
		maintainpart = &Maintaininfo[0].SCHEDULER
	default:

	}
	filePath := "./conf/maintainers.json"
	file, err := os.Open(filePath)
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return -1
	}
	defer file.Close()
	j := 0
	flag := 0
	for _, element := range maintainers[1:] {
		for _, value := range *maintainpart {
			if value != element.ID {
				(*maintainpart)[j] = value
				j++

			} else {
				flag = 1
			}
		}
	}
	*maintainpart = (*maintainpart)[:j]
	encoder := json.NewEncoder(file)
	err = encoder.Encode(Maintaininfo)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return -1
	}
	return flag
}

// 查询负责人语句format：@bot 查询负责人 端名 @负责人
// eg. @小刻-测试中 查询负责人 RUST
func searchMaintainer(data *dto.WSATMessageData) string {
	var maintainpart *[]string
	pattern := regexp.MustCompile(`查询负责人\s(\w+)\b`)
	matches := pattern.FindStringSubmatch(data.Content)
	part := ""
	if len(matches) > 1 {
		part = matches[1]
	}
	switch part {
	case "RUST":
		maintainpart = &Maintaininfo[0].RUST
	case "FETCHER":
		maintainpart = &Maintaininfo[0].FETCHER
	case "ANALYZER":
		maintainpart = &Maintaininfo[0].ANALYZER
	case "SCHEDULER":
		maintainpart = &Maintaininfo[0].SCHEDULER
	default:

	}
	content := ""
	for _, value := range *maintainpart {
		content += "<@" + value + "> "
	}
	return content
}
