package chatroom

import (
	"context"
	"sync"

	"github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/chat/pkg/websocket"
	"github.com/esmailemami/chess/shared/database/redis"
	sharedModels "github.com/esmailemami/chess/shared/models"
	sharedWebsocket "github.com/esmailemami/chess/shared/websocket"
	"github.com/google/uuid"
)

type ChatRoom struct {
	mutex       sync.Mutex
	connections map[uuid.UUID]map[uuid.UUID]*sharedWebsocket.Client

	messageService *service.MessageService
	roomService    *service.RoomService

	roomID uuid.UUID
	wss    sharedWebsocket.Server
}

func NewChatRoom(roomID uuid.UUID, wss sharedWebsocket.Server) *ChatRoom {
	var (
		redisCache     = redis.GetConnection()
		messageService = service.NewMessageService(redisCache)
		roomService    = service.NewRoomService(redisCache)
	)

	return &ChatRoom{
		mutex:          sync.Mutex{},
		connections:    make(map[uuid.UUID]map[uuid.UUID]*sharedWebsocket.Client),
		messageService: messageService,
		roomService:    roomService,
		roomID:         roomID,
		wss:            wss,
	}
}

func (g *ChatRoom) Connect(client *sharedWebsocket.Client) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.connect(client)

	room, err := g.roomService.Get(context.Background(), g.roomID)

	if err != nil {
		g.wss.SendErrorMessageToClient(client.SessionID, err.Error())
		return
	}

	lastMessages, err := g.messageService.GetLastMessages(client.Context, g.roomID)

	if err != nil {
		g.wss.SendErrorMessageToClient(client.SessionID, err.Error())
		return
	}

	g.wss.SendMessageToClient(client.SessionID, websocket.RoomDetail, &RoomMessage{
		RoomID: g.roomID,
		Data: &RoomOutPutModel{
			Room:     room,
			Messages: lastMessages,
		},
	})
}

func (g *ChatRoom) JoinUser(user *sharedModels.User) {
	g.mutex.Lock()

	data := &models.RoomUserOutPutModel{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.UserJoined, &RoomMessage{
				RoomID: g.roomID,
				Data:   data,
			})
		}
	}

	g.mutex.Unlock()
}

func (g *ChatRoom) LeftUser(user *sharedModels.User) {
	g.mutex.Lock()

	for _, client := range g.connections[user.ID] {
		g.disconnect(client)
	}

	data := &models.RoomUserOutPutModel{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.UserLeft, &RoomMessage{
				RoomID: g.roomID,
				Data:   data,
			})
		}
	}

	g.mutex.Unlock()
}

func (g *ChatRoom) Disconnect(client *sharedWebsocket.Client) {
	g.mutex.Lock()
	g.disconnect(client)
	g.mutex.Unlock()
}

func (g *ChatRoom) SendMessage(req *sharedWebsocket.ClientMessage[websocket.NewMessageRequest]) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	msg, err := g.messageService.NewMessage(req.Ctx, req.Data.RoomID, req.UserID, req.Data.Content, req.Data.ReplyTo)

	if err != nil {
		g.wss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.NewMessage, &RoomMessage{
				RoomID: g.roomID,
				Data:   msg,
			})
		}
	}
}

func (g *ChatRoom) EditMessage(req *sharedWebsocket.ClientMessage[websocket.EditMessageRequest]) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	msg, err := g.messageService.EditMessage(req.Ctx, req.Data.ID, req.Data.Content)

	if err != nil {
		g.wss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.EditMessage, &RoomMessage{
				RoomID: g.roomID,
				Data:   msg,
			})
		}
	}
}

func (g *ChatRoom) DeleteMessage(req *sharedWebsocket.ClientMessage[websocket.DeleteMessageRequest]) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	err := g.messageService.DeleteMessage(req.Ctx, req.Data.ID)

	if err != nil {
		g.wss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.DeleteMessage, &RoomMessage{
				RoomID: g.roomID,
				Data: struct {
					ID uuid.UUID `json:"id"`
				}{req.Data.ID},
			})
		}
	}
}

func (g *ChatRoom) connect(client *sharedWebsocket.Client) {
	userClients, ok := g.connections[client.UserID]
	if !ok {
		userClients = make(map[uuid.UUID]*sharedWebsocket.Client)
	}
	userClients[client.SessionID] = client
	g.connections[client.UserID] = userClients
}

func (g *ChatRoom) disconnect(client *sharedWebsocket.Client) {
	userClients, ok := g.connections[client.UserID]
	if !ok {
		return
	}

	delete(userClients, client.SessionID)

	if len(userClients) == 0 {
		delete(g.connections, client.UserID)
	} else {
		g.connections[client.UserID] = userClients
	}
}
