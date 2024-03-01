package models

import (
	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

type Message struct {
	models.Model

	Content   string     `gorm:"content" json:"content"`
	ReplyToID *uuid.UUID `gorm:"reply_to_id" json:"replyToId"`
	ReplyTo   *Message   `gorm:"foreignKey:reply_to_id;references:id;" json:"replyTo"`
	RoomID    uuid.UUID  `gorm:"room_id" json:"roomId"`
	Room      *Room      `gorm:"foreignKey:room_id;references:id;" json:"room"`
	IsEdited  bool       `gorm:"column:is_edited" json:"isEdited"`
	IsSeen    bool       `gorm:"column:is_seen" json:"isSeen"`
	Type      string     `gorm:"column:type" json:"type"`
}

func (Message) TableName() string {
	return "chat.message"
}

const (
	MESSAGE_TYPE_TEXT  = "text"
	MESSAGE_TYPE_IMAGE = "image"
	MESSAGE_TYPE_VIDEO = "video"
)
