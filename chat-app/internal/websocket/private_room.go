package websocket

import (
	"github.com/esmailemami/chess/shared/websocket"
)

var (
	PrivateRoomRegisterCh         = make(chan *websocket.Client, 256)
	PrivateRoomUnregisterCh       = make(chan *websocket.Client, 256)
	PrivateRoomNewMessageCh       = make(chan *websocket.ClientMessage[NewMessageRequest], 256)
	PrivateRoomEditMessageCh      = make(chan *websocket.ClientMessage[EditMessageRequest], 256)
	PrivateRoomDeleteMessageCh    = make(chan *websocket.ClientMessage[DeleteMessageRequest], 256)
	PrivateRoomSeenMessageCh      = make(chan *websocket.ClientMessage[SeenMessageRequest], 256)
	PrivateRoomIsTypingCh         = make(chan *websocket.ClientMessage[IsTypingRequest], 256)
	PrivateRoomPinMessageCh       = make(chan *websocket.ClientMessage[PinMessageRequest], 256)
	PrivateRoomDeletePinMessageCh = make(chan *websocket.ClientMessage[PinMessageRequest], 256)
)

func PrivateChatRoomOnMessage(c *websocket.Client, msg *websocket.Message) {
	switch msg.Type {
	case NewMessage:
		var req NewMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PrivateRoomNewMessageCh <- websocket.NewClientMessage(c, req)
	case EditMessage:
		var req EditMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PrivateRoomEditMessageCh <- websocket.NewClientMessage(c, req)
	case DeleteMessage:
		var req DeleteMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PrivateRoomDeleteMessageCh <- websocket.NewClientMessage(c, req)
	case SeenMessage:
		var req SeenMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PrivateRoomSeenMessageCh <- websocket.NewClientMessage(c, req)

	case IsTyping:
		var req IsTypingRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PrivateRoomIsTypingCh <- websocket.NewClientMessage(c, req)
	case PinMessage:
		var req PinMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PrivateRoomPinMessageCh <- websocket.NewClientMessage(c, req)
	case DeletePinMessage:
		var req PinMessageRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		PrivateRoomDeletePinMessageCh <- websocket.NewClientMessage(c, req)
	}
}

func PrivateChatRoomOnRegister(c *websocket.Client) {
	PrivateRoomRegisterCh <- c
}

func PrivateChatRoomOnUnregister(c *websocket.Client) {
	PrivateRoomUnregisterCh <- c
}
