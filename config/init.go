package config

import (
	"github.com/MikaBot-Project/MikaPluginLib/pluginConfig"
	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
)

var configReloadFunc func(msg pluginIO.Message)
var config = struct {
	WebsocketHost    string   `json:"websocket_host"`
	Commands         []string `json:"commands"`
	Message          bool     `json:"message"`
	Prefixes         []string `json:"prefixes"`
	NoPrefixCommands []string `json:"no_prefix_commands"`
}{"",
	[]string{},
	false,
	[]string{},
	[]string{},
}
var WebsocketHost string
var Commands []string
var Message bool
var Prefixes []string
var NoPrefixCommands []string

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
	Commands = config.Commands
	Message = config.Message
	Prefixes = config.Prefixes
	NoPrefixCommands = config.NoPrefixCommands
}
