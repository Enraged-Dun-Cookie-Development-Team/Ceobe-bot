package main

import (
	"ceobe-bot/bootstrap"
)

func main() {
	bootstrap.GetConfig();
	bootstrap.InitServer();
}

