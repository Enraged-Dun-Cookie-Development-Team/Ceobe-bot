package bootstrap

import (
	"ceobe-bot/initialize"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Grpc initialize.GrpcConfig `json:"grpc" envPrefix:"GRPC__"`
	Bot  initialize.BotConfig  `json:"bot" envPrefix:"BOT__"`
}
var config Config

func GetConfig() {
	if err := env.Parse(&config, env.Options{
		Prefix: "CEOBE_",
	}); err != nil {
		fmt.Printf("环境变量读取错误，尝试读取json格式变量")
	}
	fmt.Printf("%+v\n", config)
}