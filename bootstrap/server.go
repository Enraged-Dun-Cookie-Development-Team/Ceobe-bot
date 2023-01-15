package bootstrap

import "ceobe-bot/initialize"

func InitServer() {
	initialize.InitBot();
	initialize.InitGrpc();
}