package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Server interface {
	Run()
	HandleWS(ctx *gin.Context)
	BroadCastMessage(msgType string, content any) error
	SendClientToClientMessage(senderID, receiverID uuid.UUID, msgType string, content any) error
	SendMessageToClient(clientID uuid.UUID, msgType string, content any) error
	SendClientToUserMessage(clientID uuid.UUID, userID uuid.UUID, msgType string, content any) error
	SendMessageToUser(userID uuid.UUID, msgType string, content any) error
	SendErrorMessageToClient(clientID uuid.UUID, err string) error
	Shutdown()
	OnRegister(fn func(*Client))
	OnUnregister(fn func(*Client))
}

type DefaultServer struct {
	upgrader       websocket.Upgrader
	clientsMutex   sync.Mutex
	broadcast      chan *Message
	register       chan *Client
	unregister     chan *Client
	clients        map[uuid.UUID]*Client
	shutdownSignal chan struct{}
	onMessage      func(*Client, *Message)
	onRegisterFn   func(*Client)
	onUnregisterFn func(*Client)
}

func NewServer(onMessage func(*Client, *Message)) Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	return &DefaultServer{
		upgrader:       upgrader,
		broadcast:      make(chan *Message),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		shutdownSignal: make(chan struct{}),
		clients:        make(map[uuid.UUID]*Client),
		onMessage:      onMessage,
	}
}

func (s *DefaultServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.clientsMutex.Lock()
			s.clients[client.SessionID] = client

			logging.Info("Websocket client registered", "clientId", client.SessionID)

			if s.onRegisterFn != nil {
				s.onRegisterFn(client)
			}

			s.clientsMutex.Unlock()

		case client := <-s.unregister:
			s.clientsMutex.Lock()
			delete(s.clients, client.SessionID)

			logging.Info("Websocket client unregistered", "clientId", client.SessionID)

			if s.onUnregisterFn != nil {
				s.onUnregisterFn(client)
			}

			s.clientsMutex.Unlock()

		case message := <-s.broadcast:
			s.broadcastMessage(message)

		case <-s.shutdownSignal:
			logging.Info("WebSocket server shutting down...")
			s.clientsMutex.Lock()
			for _, client := range s.clients {
				client.disconnectGracefully()
			}
			s.clientsMutex.Unlock()

			<-time.After(5 * time.Second)

			return
		}
	}
}

func (s *DefaultServer) HandleWS(ctx *gin.Context) {
	conn, err := s.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(ctx, s, conn)

	s.register <- client
	defer func() {
		client.disconnectGracefully()
	}()

	logging.Info("Websocket client connected", "clientId", client.SessionID, "remoteAddr", ctx.Request.RemoteAddr)

	go client.writeMessages()
	client.readMessages()
}

func (s *DefaultServer) BroadCastMessage(msgType string, content any) error {
	msg, err := NewMessage(msgType, content, "server broadcaster")

	if err != nil {
		return err
	}

	s.broadcast <- msg

	return nil
}

func (s *DefaultServer) SendClientToClientMessage(senderID, receiverID uuid.UUID, msgType string, content any) error {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	sender, ok := s.clients[senderID]

	if !ok {
		logging.Warn("Websocket client not found", "senderId", senderID)

		return fmt.Errorf("sender (%s) not found ", senderID)
	}

	receiver, ok := s.clients[receiverID]

	if !ok {
		logging.Warn("Websocket client not found", "receiverId", receiver)

		return fmt.Errorf("receiver (%s) not found", receiverID)
	}

	msg, err := NewMessage(msgType, content, sender.User.Username)

	if err != nil {
		return err
	}

	receiver.send <- msg
	return nil
}

func (s *DefaultServer) SendMessageToClient(clientID uuid.UUID, msgType string, content any) error {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	receiver, ok := s.clients[clientID]

	if !ok {
		logging.Warn("Websocket client not found", "clientId", clientID)

		return fmt.Errorf("receiver (%s) not found", clientID)
	}

	msg, err := NewMessage(msgType, content, "Server")

	if err != nil {
		return err
	}

	receiver.send <- msg
	return nil
}
func (s *DefaultServer) SendClientToUserMessage(clientID uuid.UUID, userID uuid.UUID, msgType string, content any) error {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	sender, ok := s.clients[clientID]

	if !ok {
		logging.Warn("Websocket client not found", "senderId", clientID)
	}

	msg, err := NewMessage(msgType, content, sender.User.Username)

	if err != nil {
		return err
	}

	for _, client := range s.clients {
		if client.UserID == userID {
			client.send <- msg
		}
	}

	return nil

}
func (s *DefaultServer) SendMessageToUser(userID uuid.UUID, msgType string, content any) error {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	msg, err := NewMessage(msgType, content, "Server")

	if err != nil {
		return err
	}

	for _, client := range s.clients {
		if client.UserID == userID {
			client.send <- msg
		}
	}

	return nil
}

func (s *DefaultServer) Shutdown() {
	close(s.shutdownSignal)
}

func (s *DefaultServer) SendErrorMessageToClient(clientID uuid.UUID, err string) error {
	return s.SendMessageToClient(clientID, "error", ErrorMessage{Message: err})
}

func (s *DefaultServer) OnRegister(fn func(*Client)) {
	s.onRegisterFn = fn
}

func (s *DefaultServer) OnUnregister(fn func(*Client)) {
	s.onUnregisterFn = fn
}

func (s *DefaultServer) broadcastMessage(message *Message) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	var wg sync.WaitGroup
	ch := make(chan *Client, 5)

	go func() {
		for client := range ch {
			client.send <- message
			wg.Done()
		}
	}()

	for _, client := range s.clients {
		wg.Add(1)
		ch <- client
	}

	close(ch)
	wg.Wait()
}
