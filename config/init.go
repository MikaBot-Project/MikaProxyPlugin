package config

import (
	"github.com/MikaBot-Project/MikaPluginLib/pluginConfig"
	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
)

var configReloadFunc func(msg pluginIO.Message)
var config = struct {
	WebsocketHost string   `json:"websocket_host"`
	Command       []string `json:"command"`
	Message       bool     `json:"message"`
}{"", []string{}, false}
var WebsocketHost string
var Command []string
var Message bool

func init() {
	loadConfig()
	configReloadFunc = pluginIO.OperatorMap["config"]
	pluginIO.OperatorMap["config"] = func(message pluginIO.Message) {
		configReloadFunc(message)
		loadConfig()
	}
}

func loadConfig() {
	pluginConfig.ReadJson("config.json", &config)
	WebsocketHost = config.WebsocketHost
	Command = config.Command
	Message = config.Message
}
