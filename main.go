package main

import (
	"log"

	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
)

func main() {
	var data pluginIO.Message
	for {
		data = <-pluginIO.MessageChan
		log.Println("type: ", data.PostType)
	}
}
