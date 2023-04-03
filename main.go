package main

import (
	"github.com/FelixWuu/chatgpt_self_use/config"
	"github.com/FelixWuu/chatgpt_self_use/service"
)

func init() {
	config.Init()
}

func main() {
	service.StartWebService()
}
