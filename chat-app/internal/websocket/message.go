package websocket

import (
	"github.com/esmailemami/chess/shared/websocket"
	"github.com/google/uuid"
)

const (
	NewMessage         = "new-message"
	EditMessage        = "edit-message"
	DeleteMessage      = "delete-message"
	SeenMessage        = "seen-message"
	MessagesList       = "messages-list"
	RoomDetail         = "room-detail"
	UserJoined         = "user-joined"
	UserLeft           = "user-left"
	DeleteRoom         = "delete-room"
	RoomAvatarChanged  = "room-avatar-changed"
	UserProfileChanged = "user-profile-changed"
	EditRoom           = "edit-room"
	WatchRoom          = "watch-room"
	DeletetWatch       = "delete-watch"
	IsTyping           = "is-typing"
)

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

type SeenMessageRequest struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
}

type RoomRequest struct {
	Client *websocket.Client `json:"-"`
	RoomID uuid.UUID         `json:"roomId"`
}

type IsTypingRequest struct {
	RoomID uuid.UUID `json:"roomId"`
}
