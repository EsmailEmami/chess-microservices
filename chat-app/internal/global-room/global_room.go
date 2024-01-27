package globalroom

import (
	"sync"
	"time"

	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/chat/internal/models"
	"github.com/esmailemami/chess/chat/pkg/websocket"
	"github.com/esmailemami/chess/shared/database/redis"
	sharedWebsocket "github.com/esmailemami/chess/shared/websocket"
	"github.com/google/uuid"
)

var globalRoom *GlobalRoom

type GlobalRoom struct {
	mutex       sync.Mutex
	connections map[uuid.UUID]map[uuid.UUID]*sharedWebsocket.Client

	messageService *service.MessageService

	roomID uuid.UUID
}

func NewGlobalRoom() *GlobalRoom {
	messageService := service.NewMessageService(redis.GetConnection())

	return &GlobalRoom{
		mutex:          sync.Mutex{},
		connections:    make(map[uuid.UUID]map[uuid.UUID]*sharedWebsocket.Client),
		messageService: messageService,
		roomID:         models.GlobalRoomID,
	}
}

type message struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
}

func NewMessage(firstName, lastName *string, username, content string) *message {
	return &message{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Content:   content,
		Date:      time.Now(),
	}
}

func (g *GlobalRoom) Connect(client *sharedWebsocket.Client) {
	g.mutex.Lock()

	// register the client
	userClients, ok := g.connections[client.UserID]
	if !ok {
		userClients = make(map[uuid.UUID]*sharedWebsocket.Client)
	}
	userClients[client.SessionID] = client
	g.connections[client.UserID] = userClients

	// send last messages to client
	lastMessages, err := g.messageService.GetLastMessages(client.Context, g.roomID)

	if err != nil {
		websocket.GlobalRoomWss.SendErrorMessageToClient(client.SessionID, err.Error())
	} else {
		websocket.GlobalRoomWss.SendMessageToClient(client.SessionID, websocket.MessagesList, lastMessages)
	}

	g.mutex.Unlock()
}

func (g *GlobalRoom) Disconnect(client *sharedWebsocket.Client) {
	g.mutex.Lock()

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

	g.mutex.Unlock()
}

func (g *GlobalRoom) SendMessage(req *sharedWebsocket.ClientMessage[websocket.NewMessageRequest]) {
	g.mutex.Lock()

	msg, err := g.messageService.NewMessage(req.Ctx, g.roomID, req.UserID, &req.Data)

	if err != nil {
		websocket.GlobalRoomWss.SendErrorMessageToClient(req.ClientID, err.Error())
	} else {
		// send new message to all connected clients
		for _, clients := range g.connections {
			for _, client := range clients {
				websocket.GlobalRoomWss.SendMessageToClient(client.SessionID, websocket.NewMessage, msg)
			}
		}
	}

	g.mutex.Unlock()
}
