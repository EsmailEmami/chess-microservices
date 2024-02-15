package models

import (
	"github.com/esmailemami/chess/auth/internal/consts"
	baseconsts "github.com/esmailemami/chess/shared/consts"
	"github.com/esmailemami/chess/shared/models"
	"github.com/esmailemami/chess/shared/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInputModel struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
}

func (model RegisterInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Username,
			validation.Required.Error(baseconsts.Required),
			validation.By(validations.UserName()),
			validation.By(validations.NotExistsInDB(&models.User{}, "username", consts.UsernameAlreadyExists)),
		),
		validation.Field(
			&model.Password,
			validation.Required.Error(baseconsts.Required),
			validation.By(validations.StrongPassword()),
		),
	)
}

func (model RegisterInputModel) ToDBModel() (*models.User, error) {
	dbData := &models.User{
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Username:  model.Username,
		RoleID:    models.ROLE_ADMIN,
		Enabled:   true,
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	dbData.Password = string(pass)

	return dbData, nil
}
