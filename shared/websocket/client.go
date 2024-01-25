package websocket

import (
	"context"
	"encoding/json"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	sessionID      uuid.UUID
	userID         uuid.UUID
	user           *models.User
	conn           *websocket.Conn
	send           chan *Message
	ShutdownSignal chan struct{}
	wss            *DefaultServer
	context        context.Context
}

func NewClient(ctx *gin.Context, server *DefaultServer, conn *websocket.Conn) *Client {
	user := ctx.Value("user").(*models.User)

	c := &Client{
		sessionID:      uuid.New(),
		userID:         user.ID,
		user:           user,
		conn:           conn,
		send:           make(chan *Message),
		ShutdownSignal: make(chan struct{}),
		wss:            server,
		context:        ctx,
	}

	return c
}

func (client *Client) readMessages() {
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
loop:
	for {
		select {
		case message := <-client.send:

			err := client.conn.WriteJSON(message)

			if err != nil {
				logging.ErrorE("Websocket client unable to send message", err, "clientId", client.sessionID)
				continue
			}

		case <-client.ShutdownSignal:
			break loop
		}
	}

	close(client.send)
}

func (client *Client) disconnectGracefully() {
	close(client.ShutdownSignal)
	client.conn.Close()
	client.wss.unregister <- client
}

func (client *Client) logCloseError(code string, err error) {
	logging.WarnE("Websocket client disconnected("+code+")", err, "clientId", client.sessionID)
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
		logging.ErrorE("Websocket client disconnected unmarshalling data", err, "clientId", c.sessionID)
		return false
	}

	return true
}
