package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Type    string  `json:"type"`
	Unix    int64   `json:"unix"`
	Content string  `json:"content"`
	From    *string `json:"from,omitempty"`
}

func NewMessage(msgType string, content any, from string) (*Message, error) {

	bytes, err := json.Marshal(&content)

	if err != nil {
		return nil, err
	}

	return &Message{
		Type:    msgType,
		Unix:    time.Now().UnixNano(),
		Content: string(bytes),
		From:    &from,
	}, nil
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type NotificationMessage struct {
	Message string `json:"message"`
}

type ClientMessage[T any] struct {
	ClientID uuid.UUID
	UserID   uuid.UUID
	Ctx      context.Context
	Data     T
}

func NewClientMessage[T any](c *Client, data T) *ClientMessage[T] {
	return &ClientMessage[T]{
		ClientID: c.SessionID,
		UserID:   c.UserID,
		Ctx:      c.Context,
		Data:     data,
	}
}
