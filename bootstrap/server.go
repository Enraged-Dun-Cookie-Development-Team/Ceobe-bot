package bootstrap

import "ceobe-bot/initialize"

func InitServer() {
	go func() {
		initialize.InitBot()
	}()
	initialize.InitGrpc()
}
