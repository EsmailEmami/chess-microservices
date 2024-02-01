package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	SessionID      uuid.UUID
	UserID         uuid.UUID
	User           *models.User
	conn           *websocket.Conn
	send           chan *Message
	ShutdownSignal chan struct{}
	wss            *DefaultServer
	Context        context.Context
}

func NewClient(ctx *gin.Context, server *DefaultServer, conn *websocket.Conn) *Client {
	user := ctx.Value("user").(*models.User)

	c := &Client{
		SessionID:      uuid.New(),
		UserID:         user.ID,
		User:           user,
		conn:           conn,
		send:           make(chan *Message),
		ShutdownSignal: make(chan struct{}),
		wss:            server,
		Context:        ctx,
	}

	return c
}

func (client *Client) readMessages() {
	// read deadline (ping pong)
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {

		if client.wss.onPongFn != nil {
			client.wss.onPongFn(client)
		}

		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

loop:
	for {
		var msg Message
		err := client.conn.ReadJSON(&msg)

		if err != nil {
			client.handleReadMessageErr(err)
			break loop
		}

		client.wss.onMessage(client, &msg)
	}
}

func (client *Client) writeMessages() {

	ticker := time.NewTicker(pingPeriod)

loop:
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				break loop
			}

			err := client.conn.WriteJSON(message)

			if err != nil {
				logging.ErrorE("Websocket client unable to send message", err, "clientId", client.SessionID)
				continue
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-client.ShutdownSignal:
			break loop
		}
	}

	close(client.send)
	ticker.Stop()
}

func (client *Client) disconnectGracefully() {
	close(client.ShutdownSignal)
	client.conn.Close()
	client.wss.unregister <- client
}

func (client *Client) logCloseError(code string, err error) {
	logging.WarnE("Websocket client disconnected("+code+")", err, "clientId", client.SessionID)
}

func (client *Client) handleReadMessageErr(err error) {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		client.logCloseError("CloseNormalClosure", err)
	} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
		client.logCloseError("CloseGoingAway", err)
	} else if websocket.IsCloseError(err, websocket.CloseProtocolError) {
		client.logCloseError("CloseProtocolError", err)
	} else if websocket.IsCloseError(err, websocket.CloseUnsupportedData) {
		client.logCloseError("CloseUnsupportedData", err)
	} else if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
		client.logCloseError("CloseNoStatusReceived", err)
	} else if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
		client.logCloseError("CloseAbnormalClosure", err)
	} else if websocket.IsCloseError(err, websocket.CloseInvalidFramePayloadData) {
		client.logCloseError("CloseInvalidFramePayloadData", err)
	} else if websocket.IsCloseError(err, websocket.ClosePolicyViolation) {
		client.logCloseError("ClosePolicyViolation", err)
	} else if websocket.IsCloseError(err, websocket.CloseMessageTooBig) {
		client.logCloseError("CloseMessageTooBig", err)
	} else if websocket.IsCloseError(err, websocket.CloseMandatoryExtension) {
		client.logCloseError("CloseMandatoryExtension", err)
	} else if websocket.IsCloseError(err, websocket.CloseInternalServerErr) {
		client.logCloseError("CloseInternalServerError", err)
	} else if websocket.IsCloseError(err, websocket.CloseServiceRestart) {
		client.logCloseError("CloseServiceRestart", err)
	} else if websocket.IsCloseError(err, websocket.CloseTryAgainLater) {
		client.logCloseError("CloseTryAgainLater", err)
	} else if websocket.IsCloseError(err, websocket.CloseTLSHandshake) {
		client.logCloseError("CloseTLSHandshake", err)
	} else {
		// Handle other close errors or unexpected errors
		client.logCloseError("UnexpectedError", err)
	}
}

func (c *Client) Unmarshal(content string, model any) bool {
	if err := json.Unmarshal([]byte(content), model); err != nil {
		logging.ErrorE("Websocket client disconnected unmarshalling data", err, "clientId", c.SessionID)
		return false
	}

	return true
}
