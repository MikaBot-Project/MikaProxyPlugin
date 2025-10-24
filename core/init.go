package core

import (
	"MikaProxyPlugin/config"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
	"github.com/lxzan/gws"
)

var sendChan chan []byte
var selfId = int64(0)

func init() {
	pluginIO.MessageRegister(func(message pluginIO.Message) {
		if config.Message {
			getMessage(message)
		}
	})
	for _, name := range config.Command {
		pluginIO.CommandRegister(name, getMessage)
	}
	pluginIO.OperatorMap["return"] = func(message pluginIO.Message) {
		if message.SubType == "self_id" {
			selfId, _ = strconv.ParseInt(message.CommandArgs[0], 10, 64)
			log.Println("selfId:", selfId)
			go Start()
		}
	}
	pluginIO.SendOperator("panel", "get_self_id", []string{})
	sendChan = make(chan []byte, 10)
}

type logger struct{}

func (l logger) Error(v ...interface{}) {
	log.Println(v...)
}

func Start() {
	header := http.Header{}
	header.Set("X-Self-ID", strconv.FormatInt(selfId, 10))
	header.Set("User-Agent", "OneBot")
	header.Set("x-client-role", "Universal")
	header.Set("Authorization", "Bearer "+pluginIO.RandomString(16))
	conn, _, err := gws.NewClient(&SocketHandler{}, &gws.ClientOption{
		ParallelEnabled:     true,
		Recovery:            gws.Recovery,
		Addr:                config.WebsocketHost,
		Logger:              logger{},
		RequestHeader:       header,
		WriteMaxPayloadSize: 1024 * 1024,
	})
	if err != nil {
		log.Println(err)
	}
	conn.ReadLoop()
}

func getMessage(msg pluginIO.Message) {
	if msg.PostType == "command" {
		msg.PostType = "message"
	}
	data, _ := json.Marshal(msg)
	sendChan <- data
}
