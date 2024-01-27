package globalroom

import (
	"github.com/esmailemami/chess/chat/pkg/websocket"
)

func Run() {
	globalRoom = NewGlobalRoom()

	for {
		select {
		case client := <-websocket.GlobalRoomRegisterCh:
			globalRoom.Connect(client)

		case client := <-websocket.GlobalRoomUnregisterCh:
			globalRoom.Disconnect(client)

		case req := <-websocket.GlobalRoomNewMessageCh:
			globalRoom.SendMessage(req)
		}
	}
}
