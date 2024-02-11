package models

import (
	"time"

	"github.com/esmailemami/chess/shared/models"
)

type User struct {
	models.User

	LastConnection time.Time `gorm:"column:last_connection" json:"lastConnection"`
	Profile        string    `gorm:"column:profile" json:"profile"`
}
