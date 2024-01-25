package models

import "github.com/esmailemami/chess/shared/models"

type Room struct {
	models.Model

	Name      string     `gorm:"name" json:"name"`
	IsPrivate bool       `gorm:"is_private" json:"isPrivate"`
	Users     []UserRoom `gorm:"users" json:"users"`
}
