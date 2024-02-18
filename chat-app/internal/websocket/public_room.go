package websocket

import (
	"github.com/esmailemami/chess/shared/websocket"
)

var (
	PublicRoomRegisterCh      = make(chan *websocket.Client, 256)
	PublicRoomUnregisterCh    = make(chan *websocket.Client, 256)
	PublicRoomNewMessageCh    = make(chan *websocket.ClientMessage[NewMessageRequest], 256)
	PublicRoomEditMessageCh   = make(chan *websocket.ClientMessage[EditMessageRequest], 256)
	PublicRoomDeleteMessageCh = make(chan *websocket.ClientMessage[DeleteMessageRequest], 256)
	PublicRoomWatchCh         = make(chan *RoomRequest, 256)
	PublicRoomDeleteWatchCh   = make(chan *RoomRequest, 256)
)

func PublicChatRoomOnMessage(c *websocket.Client, msg *websocket.Message) {
	switch msg.Type {
	case NewMessage:
		var req NewMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PublicRoomNewMessageCh <- websocket.NewClientMessage(c, req)
	case EditMessage:
		var req EditMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PublicRoomEditMessageCh <- websocket.NewClientMessage(c, req)
	case DeleteMessage:
		var req DeleteMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PublicRoomDeleteMessageCh <- websocket.NewClientMessage(c, req)
	case WatchRoom:
		var req RoomRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		req.Client = c

		PublicRoomWatchCh <- &req
	case DeletetWatch:
		var req RoomRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		req.Client = c

		PublicRoomDeleteWatchCh <- &req
	}
}

func PublicChatRoomOnRegister(c *websocket.Client) {
	PublicRoomRegisterCh <- c
}

func PublicChatRoomOnUnregister(c *websocket.Client) {
	PublicRoomUnregisterCh <- c
}
