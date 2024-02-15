package websocket

import "github.com/google/uuid"

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
