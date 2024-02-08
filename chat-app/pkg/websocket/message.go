package websocket

import "github.com/google/uuid"

type NewMessageRequest struct {
	RoomID  uuid.UUID  `json:"roomId,omitempty"`
	Content string     `json:"content"`
	ReplyTo *uuid.UUID `json:"replyTo,omitempty"`
}

type EditMessageRequest struct {
	ID      uuid.UUID `json:"id"`
	RoomID  uuid.UUID `json:"roomId"`
	Content string    `json:"content"`
}

type DeleteMessageRequest struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
}
