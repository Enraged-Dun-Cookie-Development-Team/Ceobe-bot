package bootstrap

import "ceobe-bot/initialize"

func InitServer() {
	go func() {
		initialize.InitBot(config.Bot)
	}()
	initialize.InitGrpc(config.Grpc)
}
