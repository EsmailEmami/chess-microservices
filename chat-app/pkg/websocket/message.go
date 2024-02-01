package websocket

import "github.com/google/uuid"

type NewMessageRequest struct {
	RoomID  uuid.UUID  `json:"roomId,omitempty"`
	Content string     `json:"content"`
	ReplyTo *uuid.UUID `json:"replyTo,omitempty"`
}
