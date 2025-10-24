package main

import (
	"MikaProxyPlugin/config"
	_ "MikaProxyPlugin/core"
	"log"

	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
)

func main() {
	if config.WebsocketHost == "" {
		log.Println("Websocket host is required")
		return
	}
	var data pluginIO.Message
	for {
		data = <-pluginIO.MessageChan
		log.Println("type: ", data.PostType)
	}
}
