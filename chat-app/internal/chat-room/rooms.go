package chatroom

import (
	"github.com/esmailemami/chess/chat/internal/models"
	"github.com/esmailemami/chess/chat/internal/websocket"
	sharedModels "github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

var (
	globalRoom   *ChatRoom
	publicRooms  map[uuid.UUID]*ChatRoom
	privateRooms map[uuid.UUID]*ChatRoom
)

// in 'init' we do not have redis cache yet!
func initializeRooms() {
	publicRooms = make(map[uuid.UUID]*ChatRoom)
	privateRooms = make(map[uuid.UUID]*ChatRoom)
	globalRoom = NewChatRoom(models.GlobalRoomID, true, websocket.GlobalRoomWss)
}

func getPublicChatRoom(id uuid.UUID) *ChatRoom {
	room, ok := publicRooms[id]
	if !ok {
		room = NewChatRoom(id, true, websocket.PublicChatRoomWss)
		publicRooms[id] = room
	}

	return room
}

func getPrivateChatRoom(id uuid.UUID) *ChatRoom {
	room, ok := privateRooms[id]
	if !ok {
		room = NewChatRoom(id, false, websocket.PrivateChatRoomWss)
		privateRooms[id] = room
	}

	return room
}

// public rooms

func ConnectPublicRoom(roomID, userID uuid.UUID) {
	room := getPublicChatRoom(roomID)

	for _, client := range websocket.PublicChatRoomWss.GetUserConnections(userID) {
		room.Connect(client)
	}
}

func JoinRoom(roomID uuid.UUID, user *sharedModels.User) {
	room := getPublicChatRoom(roomID)

	for _, client := range websocket.PublicChatRoomWss.GetUserConnections(user.ID) {
		room.Connect(client)
	}

	room.JoinUser(user)
}

func LeftRoom(roomID uuid.UUID, user *sharedModels.User) {
	room := getPublicChatRoom(roomID)
	room.LeftUser(user)
}

// private rooms

func ConnectPrvateRoom(roomID, userID uuid.UUID) {
	room := getPrivateChatRoom(roomID)

	for _, client := range websocket.PrivateChatRoomWss.GetUserConnections(userID) {
		room.Connect(client)
	}
}

func DeleteRoom(roomID uuid.UUID) {
	if room, ok := privateRooms[roomID]; ok {
		room.Delete()
		delete(privateRooms, roomID)

		return
	}

	if room, ok := publicRooms[roomID]; ok {
		room.Delete()
		delete(publicRooms, roomID)

		return
	}
}

func RoomEdited(roomID uuid.UUID) {
	if room, ok := publicRooms[roomID]; ok {
		room.Edit()
	}
}
