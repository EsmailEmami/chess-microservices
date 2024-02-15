package models

import (
	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

type Friend struct {
	models.BaseModel

	UserID uuid.UUID `gorm:"column:user_id" json:"userId"`
	User   *User     `gorm:"foreignKey:user_id;references:id;" json:"user"`

	FriendID uuid.UUID `gorm:"column:friend_id" json:"friendId"`
	Friend   *User     `gorm:"foreignKey:friend_id;references:id;" json:"friend"`
}

func (Friend) TableName() string {
	return "public.friend"
}
