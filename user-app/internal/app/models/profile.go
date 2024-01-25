package models

import (
	"github.com/esmailemami/chess/shared/models"
	sharedValidations "github.com/esmailemami/chess/shared/validations"
	"github.com/esmailemami/chess/user/pkg/consts"
	"github.com/esmailemami/chess/user/pkg/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type UserProfileOutPutModel struct {
	ID        uuid.UUID `gorm:"id" json:"id"`
	FirstName *string   `gorm:"first_name" json:"firstName"`
	LastName  *string   `gorm:"last_name" json:"lastName"`
	Mobile    *string   `gorm:"mobile" json:"mobile"`
	Username  string    `gorm:"username" json:"username"`
	RoleID    uuid.UUID `gorm:"role_id" json:"roleId"`
	RoleName  string    `gorm:"role_name" json:"roleName"`
}

type UserChangeProfileInputModel struct {
	ID        uuid.UUID `json:"-"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Mobile    *string   `json:"mobile"`
	Username  string    `json:"username"`
}

func (model UserChangeProfileInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Mobile,
			validation.By(validations.IsValidMobileNumber()),
		),
		validation.Field(
			&model.Username,
			validation.By(sharedValidations.UserName()),
			validation.By(sharedValidations.NotExistsInDBWithCond(&models.User{}, "username", consts.UsernameAlreadyExists, "id != ?", model.ID)),
		),
	)
}

func (model UserChangeProfileInputModel) MergeDBModel(dbData *models.User) {
	dbData.FirstName = model.FirstName
	dbData.LastName = model.FirstName
	dbData.Username = model.Username
	dbData.Mobile = model.Mobile
}
