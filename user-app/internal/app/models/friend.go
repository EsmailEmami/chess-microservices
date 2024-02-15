package models

import "github.com/google/uuid"

type FriendQueryParams struct {
	SearchTerm string `json:"searchTerm"`
}

type FriendOutPutModel struct {
	ID        uuid.UUID `gorm:"column:id" json:"id"`
	FirstName string    `gorm:"column:first_name" json:"firstName"`
	LastName  string    `gorm:"column:last_name" json:"lastName"`
	Username  string    `gorm:"column:username" json:"username"`
	Profile   string    `gorm:"column:profile" json:"profile"`
}
