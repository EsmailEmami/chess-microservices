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
}

func (Message) TableName() string {
	return "chat.message"
}
