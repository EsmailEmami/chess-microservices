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
}
