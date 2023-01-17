package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Grpc GrpcConfig `json:"grpc" envPrefix:"GRPC__"`
	Bot  BotConfig  `json:"bot" envPrefix:"BOT__"`
}

var config Config

func SetConfig() {
	if err := env.Parse(&config, env.Options{
		Prefix: "CEOBE_",
	}); err != nil {
		fmt.Println("环境变量读取错误，尝试读取json格式变量")
		path, _ := os.Getwd()
		// 打开文件
		file, _ := os.Open(path + "/conf/config.json")
		// 关闭文件
		defer file.Close()
		// NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
		decoder := json.NewDecoder(file)
		//Decode从输入流读取下一个json编码值并保存在v指向的值里
		err = decoder.Decode(&config)
		if err != nil {
			log.Fatal(err)
		}
		validate := validator.New()
		err := validate.Struct(config)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetConfig() Config {
	return config
}
