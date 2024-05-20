package conf

type BotConfig struct {
	AppId         int64  `json:"app_id" env:"APP_ID,notEmpty" validate:"required"`
	Token         string `json:"token" env:"TOKEN,notEmpty" validate:"required"`
	ChannelNotice string `json:"channel_notice" env:"CHANNEL_NOTICE,notEmpty" validate:"required"`
	ChannelInfo string `json:"channel_info" env:"CHANNEL_INFO"`
}
