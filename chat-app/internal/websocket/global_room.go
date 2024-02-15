package websocket

import (
	"github.com/esmailemami/chess/shared/websocket"
)

var (
	GlobalRoomRegisterCh   = make(chan *websocket.Client, 256)
	GlobalRoomUnregisterCh = make(chan *websocket.Client, 256)
	GlobalRoomNewMessageCh = make(chan *websocket.ClientMessage[NewMessageRequest], 256)
)

func GlobalRoomOnMessage(c *websocket.Client, msg *websocket.Message) {
	switch msg.Type {
	case NewMessage:
		var req NewMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		GlobalRoomNewMessageCh <- websocket.NewClientMessage(c, req)
	}
}

func GlobalRoomOnRegister(c *websocket.Client) {
	GlobalRoomRegisterCh <- c
}

func GlobalRoomOnUnregister(c *websocket.Client) {
	GlobalRoomUnregisterCh <- c
}
