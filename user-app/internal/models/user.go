package models

import (
	"time"

	"github.com/esmailemami/chess/shared/models"
)

type User struct {
	models.User

	LastConnection time.Time `gorm:"column:last_connection" json:"lastConnection"`
	Bio            string    `gorm:"column:bio" json:"bio"`

	Friends []Friend `gorm:"foreignKey:user_id;references:id;" json:"friends"`
}

func (User) TableName() string {
	return "public.user"
}
