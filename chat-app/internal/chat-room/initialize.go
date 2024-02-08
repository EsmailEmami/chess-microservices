package chatroom

import (
	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/chat/pkg/websocket"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/logging"
)

func Run() {
	initializeRooms()

	go runGlobalChatRoom()
	go runPublicChatRoom()
	go runPrivateChatRoom()
}

func runGlobalChatRoom() {
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

func runPublicChatRoom() {
	roomService := service.NewRoomService(redis.GetConnection())

	for {
		select {
		case client := <-websocket.PublicRoomRegisterCh:
			userRooms, err := roomService.GetUserRoomIDs(client.Context, client.UserID, false)
			if err != nil {
				logging.ErrorE("failed to get user rooms", err)
			}

			logging.Info("user rooms", "len", len(userRooms), "clientId", client.SessionID)

			for _, roomID := range userRooms {
				getPublicChatRoom(roomID).Connect(client)
			}

		case client := <-websocket.PublicRoomUnregisterCh:
			userRooms, err := roomService.GetUserRoomIDs(client.Context, client.UserID, false)
			if err != nil {
				logging.ErrorE("failed to get user rooms", err)
			}

			for _, roomID := range userRooms {
				getPublicChatRoom(roomID).Disconnect(client)
			}

		case req := <-websocket.PublicRoomNewMessageCh:
			room, ok := publicRooms[req.Data.RoomID]
			if ok {
				room.SendMessage(req)
			}
		case req := <-websocket.PublicRoomEditMessageCh:
			room, ok := publicRooms[req.Data.RoomID]
			if ok {
				room.EditMessage(req)
			}
		case req := <-websocket.PublicRoomDeleteMessageCh:
			room, ok := publicRooms[req.Data.RoomID]
			if ok {
				room.DeleteMessage(req)
			}
		}
	}
}

func runPrivateChatRoom() {
	roomService := service.NewRoomService(redis.GetConnection())

	for {
		select {
		case client := <-websocket.PrivateRoomRegisterCh:
			userRooms, err := roomService.GetUserRoomIDs(client.Context, client.UserID, false)
			if err != nil {
				logging.ErrorE("failed to get user rooms", err)
			}

			for _, roomID := range userRooms {
				getPrivateChatRoom(roomID).Connect(client)
			}

		case client := <-websocket.PrivateRoomUnregisterCh:
			userRooms, err := roomService.GetUserRoomIDs(client.Context, client.UserID, false)
			if err != nil {
				logging.ErrorE("failed to get user rooms", err)
			}

			for _, roomID := range userRooms {
				getPrivateChatRoom(roomID).Disconnect(client)
			}

		case req := <-websocket.PrivateRoomNewMessageCh:
			room, ok := privateRooms[req.Data.RoomID]
			if ok {
				room.SendMessage(req)
			}
		case req := <-websocket.PrivateRoomEditMessageCh:
			room, ok := privateRooms[req.Data.RoomID]
			if ok {
				room.EditMessage(req)
			}
		case req := <-websocket.PrivateRoomDeleteMessageCh:
			room, ok := privateRooms[req.Data.RoomID]
			if ok {
				room.DeleteMessage(req)
			}
		}
	}
}
