package models

import (
	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

type UserRoom struct {
	models.BaseModel

	UserID uuid.UUID    `gorm:"user_id" json:"userId"`
	User   *models.User `gorm:"foreignKey:user_id;references:id;" json:"user"`
	RoomID uuid.UUID    `gorm:"room_id" json:"roomId"`
	Room   *Room        `gorm:"foreignKey:room_id;references:id;" json:"room"`
}

func (UserRoom) TableName() string {
	return "chat.user_room"
}
