package models

import (
	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

var GlobalRoomID = uuid.MustParse("9b7af5f3-cc90-4127-96ca-1e8b32e8bb75")

type Room struct {
	models.Model

	Name      string     `gorm:"name" json:"name"`
	IsPrivate bool       `gorm:"is_private" json:"isPrivate"`
	Users     []UserRoom `gorm:"foreignKey:room_id;references:id;" json:"users"`
	Messages  []Message  `gorm:"foreignKey:room_id;references:id;" json:"messages"`
}

func (Room) TableName() string {
	return "chat.room"
}
