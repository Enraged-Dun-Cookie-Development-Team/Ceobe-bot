package main

import (
	"ceobe-bot/bootstrap"
	"ceobe-bot/conf"
)

func main() {
	conf.SetConfig()
	bootstrap.InitServer()
}
