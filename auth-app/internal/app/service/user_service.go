package service

import (
	"context"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/models"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	db := psql.DBContext(ctx).Model(&models.User{})

	var user models.User

	if err := db.Find(&user, "username=?", username).Error; err != nil {
		return nil, errs.NotFoundErr().WithError(err)
	}

	return &user, nil
}
