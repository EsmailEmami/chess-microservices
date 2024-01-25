package models

import (
	"github.com/esmailemami/chess/shared/consts"
	"github.com/esmailemami/chess/shared/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UserChangePasswordInputModel struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func (model UserChangePasswordInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Password,
			validation.Required.Error(consts.Required),
		),
		validation.Field(
			&model.NewPassword,
			validation.Required.Error(consts.Required),
			validation.By(validations.StrongPassword()),
		),
	)
}
