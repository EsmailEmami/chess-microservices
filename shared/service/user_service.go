package service

import "github.com/esmailemami/chess/shared/models"

type UserService struct {
	BaseService[models.User]
}

func NewUserService() *UserService {
	return new(UserService)
}
