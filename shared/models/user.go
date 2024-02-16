package models

import "github.com/google/uuid"

type User struct {
	Model

	FirstName *string   `gorm:"first_name" json:"firstName"`
	LastName  *string   `gorm:"last_name" json:"lastName"`
	Mobile    *string   `gorm:"mobile" json:"mobile"`
	Username  string    `gorm:"username" json:"username"`
	Password  string    `gorm:"password" json:"-"`
	RoleID    uuid.UUID `gorm:"role_id" json:"roleId"`
	Role      *Role     `gorm:"foreignKey:role_id;references:id" json:"role"`
	Enabled   bool      `gorm:"enabled" json:"enabled"`
	Profile   string    `gorm:"column:profile" json:"profile"`
}

func (User) TableName() string {
	return "public.user"
}

func (u User) IsAdmin() bool {
	if u.IsRoot() {
		return true
	}
	return u.RoleID == ROLE_ADMIN
}

func (u User) IsRoot() bool {
	return u.RoleID == ROLE_ROOT
}
