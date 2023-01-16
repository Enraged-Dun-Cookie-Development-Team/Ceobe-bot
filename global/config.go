package global

type Config struct {
	Grpc *Grpc `json:"grpc"`
	Bot  *Bot  `json:"bot"`
}

type Grpc struct {
	Port int `json:"port"`
}

type Bot struct {
	ChannelNotice string `json:"channel_notice"`
}
