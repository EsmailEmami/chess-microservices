package chatroom

import (
	"context"
	"sync"

	"github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/chat/internal/rabbitmq"
	"github.com/esmailemami/chess/chat/internal/util"
	"github.com/esmailemami/chess/chat/internal/websocket"

	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/logging"
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

	isPublic bool
}

func NewChatRoom(roomID uuid.UUID, isPublic bool, wss sharedWebsocket.Server) *ChatRoom {
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
		isPublic:       isPublic,
	}
}

func (g *ChatRoom) Connect(client *sharedWebsocket.Client) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.connect(client)

	room, err := g.roomService.Get(context.Background(), g.roomID, &client.UserID)

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
		Profile:   util.FilePathPrefix(user.Profile),
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

	// delete the room for the user
	for _, client := range g.connections[user.ID] {
		g.wss.SendMessageToClient(client.SessionID, websocket.DeleteRoom, &RoomMessage{
			RoomID: g.roomID,
		})
	}

	// delete the user from connections
	delete(g.connections, user.ID)

	// send the user left message to other connections
	data := &models.RoomUserOutPutModel{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Profile:   util.FilePathPrefix(user.Profile),
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

	msg, err := g.messageService.EditMessage(req.Ctx, req.Data.ID, req.Data.RoomID, req.Data.Content)

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

	err := g.messageService.DeleteMessage(req.Ctx, req.Data.ID, req.Data.RoomID)

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

func (g *ChatRoom) SeenMessage(req *sharedWebsocket.ClientMessage[websocket.SeenMessageRequest]) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	err := g.messageService.SeenMessage(req.Ctx, req.Data.ID, req.Data.RoomID)

	if err != nil {
		g.wss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.SeenMessage, &RoomMessage{
				RoomID: g.roomID,
				Data: struct {
					ID uuid.UUID `json:"id"`
				}{req.Data.ID},
			})
		}
	}
}

func (g *ChatRoom) Delete() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	logging.Debug("Delete room called from websocket")

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.DeleteRoom, &RoomMessage{
				RoomID: g.roomID,
			})

			g.disconnect(client)
		}
	}
}

func (g *ChatRoom) Edit() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// this is a public room always
	room, err := g.roomService.Get(context.Background(), g.roomID, nil)

	if err != nil {
		return
	}

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.EditRoom, &RoomMessage{
				RoomID: g.roomID,
				Data:   room,
			})
		}
	}
}

func (g *ChatRoom) AvatarChanged(avatar string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.RoomAvatarChanged, &RoomMessage{
				RoomID: g.roomID,
				Data: &RoomAvatarChangedModel{
					Avatar: avatar,
				},
			})
		}
	}
}

func (g *ChatRoom) UserProfileChanged(userID uuid.UUID, profile string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// send the profile changed to all users
	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.UserProfileChanged, &RoomMessage{
				RoomID: g.roomID,
				Data: &UserProfileChangedModel{
					UserID:  userID,
					Profile: profile,
				},
			})
		}
	}

	// if the room is private and both users are connected, the room avatar must be change too
	if !g.isPublic && len(g.connections) == 2 {
		for id, clients := range g.connections {
			if id == userID {
				continue
			}

			// send avatar changed message
			for _, client := range clients {
				g.wss.SendMessageToClient(client.SessionID, websocket.RoomAvatarChanged, &RoomMessage{
					RoomID: g.roomID,
					Data: &RoomAvatarChangedModel{
						Avatar: profile,
					},
				})
			}
		}
	}
}

func (g *ChatRoom) Watch(client *sharedWebsocket.Client) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if !g.isPublic {
		g.wss.SendErrorMessageToClient(client.SessionID, "this room is not public")
		return
	}

	if _, ok := g.connections[client.UserID]; ok {
		g.wss.SendErrorMessageToClient(client.SessionID, "you already joined the room")
		return
	}

	g.roomService.SetWatchRoom(client.SessionID, g.roomID)

	g.connect(client)

	room, err := g.roomService.Get(context.Background(), g.roomID, &client.UserID)

	if err != nil {
		g.wss.SendErrorMessageToClient(client.SessionID, err.Error())
		return
	}

	lastMessages, err := g.messageService.GetLastMessages(client.Context, g.roomID)

	if err != nil {
		g.wss.SendErrorMessageToClient(client.SessionID, err.Error())
		return
	}

	g.wss.SendMessageToClient(client.SessionID, websocket.WatchRoom, &RoomMessage{
		RoomID: g.roomID,
		Data: &RoomOutPutModel{
			Room:     room,
			Messages: lastMessages,
		},
	})
}

func (g *ChatRoom) DeleteWatch(client *sharedWebsocket.Client) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if !g.isPublic {
		g.wss.SendErrorMessageToClient(client.SessionID, "this room is not public")
		return
	}

	g.roomService.DeleteWatchRoom(client.SessionID, g.roomID)
	g.disconnect(client)
}

func (g *ChatRoom) IsTyping(req *sharedWebsocket.ClientMessage[websocket.IsTypingRequest]) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	clients, ok := g.connections[req.UserID]

	if !ok {
		return
	}

	requestedUser := clients[req.ClientID].User

	for _, clients := range g.connections {
		for _, client := range clients {
			g.wss.SendMessageToClient(client.SessionID, websocket.IsTyping, &RoomMessage{
				RoomID: g.roomID,
				Data: UserIsTypingModel{
					ID:        requestedUser.ID,
					FirstName: requestedUser.FirstName,
					LastName:  requestedUser.LastName,
					Username:  requestedUser.Username,
				},
			})
		}
	}
}

func (g *ChatRoom) SendFileMessage(req *rabbitmq.RoomFileMessage) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	msg, err := g.messageService.NewFileMessage(context.Background(), req.RoomID, req.UserID, req.MessageID, req.File, req.Type)

	if err != nil {
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

	// delete the room from rooms map
	if len(g.connections) == 0 {
		if g.isPublic {
			delete(publicRooms, g.roomID)
		} else {
			delete(privateRooms, g.roomID)
		}
	}
}
