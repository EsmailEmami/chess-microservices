package models

import (
	"time"

	"github.com/esmailemami/chess/shared/consts"
	"github.com/esmailemami/chess/shared/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type LoginInputModel struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IP        string `json:"-"`
	UserAgent string `json:"-"`
}

func (model LoginInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Username,
			validation.Required.Error(consts.Required),
			validation.By(validations.UserName()),
		),
		validation.Field(
			&model.Password,
			validation.Required.Error(consts.Required),
		),
	)
}

type LoginOutputModel struct {
	TokenID   uuid.UUID            `json:"-"`
	Token     string               `json:"token"`
	ExpiresAt time.Time            `json:"expiresAt"`
	ExpiresIn int64                `json:"expiresIn"`
	User      LoginOutputUserModel `json:"user"`
}

type LoginOutputUserModel struct {
	Username  string  `json:"username"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}
