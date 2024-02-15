package models

import (
	sharedValidations "github.com/esmailemami/chess/shared/validations"
	"github.com/esmailemami/chess/user/internal/app/validations"
	"github.com/esmailemami/chess/user/internal/consts"
	"github.com/esmailemami/chess/user/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type UserProfileOutPutModel struct {
	ID        uuid.UUID `gorm:"column:id" json:"id"`
	FirstName *string   `gorm:"column:first_name" json:"firstName"`
	LastName  *string   `gorm:"column:last_name" json:"lastName"`
	Mobile    *string   `gorm:"column:mobile" json:"mobile"`
	Username  string    `gorm:"column:username" json:"username"`
	RoleID    uuid.UUID `gorm:"column:role_id" json:"roleId"`
	RoleName  string    `gorm:"column:role_name" json:"roleName"`
	Profile   string    `gorm:"column:profile" json:"profile"`
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
