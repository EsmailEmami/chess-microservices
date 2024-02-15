package models

import (
	"time"

	"github.com/esmailemami/chess/shared/models"
)

type User struct {
	models.User

	LastConnection time.Time `gorm:"column:last_connection" json:"lastConnection"`
	Profile        string    `gorm:"column:profile" json:"profile"`

	Friends []Friend `gorm:"foreignKey:user_id;references:id;" json:"friends"`
}

func (User) TableName() string {
	return "public.user"
}
