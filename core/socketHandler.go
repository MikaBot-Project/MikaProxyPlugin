package core

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
	"github.com/lxzan/gws"
)

const (
	PingInterval = 10 * time.Second
	PingWait     = 100 * time.Second
)

type SocketHandler struct {
}

type metaEvent struct {
	Time          int64  `json:"time"`
	SelfId        int64  `json:"self_id"`
	PostType      string `json:"post_type"`
	MetaEventType string `json:"meta_event_type"`
	SubType       string `json:"sub_type"`
	Status        struct {
		Online bool `json:"online"`
		Good   bool `json:"good"`
	} `json:"status"`
	Interval int `json:"interval"`
}

func (c *SocketHandler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	meta := metaEvent{
		PostType:      "meta_event",
		MetaEventType: "lifecycle",
		SubType:       "connect",
		Interval:      int(PingInterval / time.Millisecond),
		SelfId:        selfId,
		Time:          time.Now().Unix(),
	}
	meta.Status.Online = true
	meta.Status.Good = true
	send, _ := json.Marshal(meta)
	_ = socket.WriteMessage(gws.OpcodeText, send)
	go func() { //发送数据
		var data []byte
		for {
			data = <-sendChan
			err := socket.WriteMessage(gws.OpcodeText, data)
			if err != nil {
				log.Println("Send data err:", err)
				sendChan <- data
				_ = socket.WriteClose(1000, nil)
				return
			}
		}
	}()
	go func() {
		meta.MetaEventType = "heartbeat"
		for {
			time.Sleep(PingInterval)
			meta.Time = time.Now().Unix()
			send, _ = json.Marshal(meta)
			err := socket.WriteMessage(gws.OpcodeText, send)
			if err != nil {
				log.Println("Send ping err:", err)
				return
			}
			_ = socket.SetDeadline(time.Now().Add(PingWait))
		}
	}()
}

func (c *SocketHandler) OnClose(socket *gws.Conn, err error) {
	log.Println("websocket close with err:", err)
	time.Sleep(PingInterval)
	go Start()
}

func (c *SocketHandler) OnPing(socket *gws.Conn, payload []byte) {
	log.Println("websocket ping")
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(payload)
}

func (c *SocketHandler) OnPong(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (c *SocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer func(message *gws.Message) {
		err := message.Close()
		if err != nil {
			log.Println(err)
		}
	}(message)
	var send = struct {
		Action string `json:"action"`
		Params any    `json:"params"`
		Echo   any    `json:"echo"`
	}{}
	var err error
	var data []byte
	switch message.Opcode {
	case gws.OpcodeText:
		err = json.Unmarshal(message.Bytes(), &send)
		if err != nil {
			log.Println("Unmarshal msg err", err)
			return
		}
		data, err = json.Marshal(send.Params)
		if err != nil {
			log.Println("Marshal msg err", err)
			return
		}
		echo, _ := json.Marshal(send.Echo)
		data = pluginIO.SendApiEcho(send.Action, data, echo)
		sendChan <- data
	case gws.OpcodePing:
		c.OnPing(socket, message.Bytes())
	}
}
