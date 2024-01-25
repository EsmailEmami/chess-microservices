package handler

import (
	"github.com/esmailemami/chess/shared/handler"
	"github.com/esmailemami/chess/user/internal/app/models"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	handler.Handler

	userService *service.UserService
}

func NewProfileHandler(userService *service.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

func (u *ProfileHandler) Profile(c *gin.Context) (*handler.Response[models.UserProfileOutPutModel], error) {
	user, err := u.GetUser(c)
	if err != nil {
		return nil, err
	}

	profile, err := u.userService.GetProfile(c, user.ID)

	if err != nil {
		return nil, err
	}

	return handler.OK(profile), nil
}

func (u *ProfileHandler) ChangePassword(c *gin.Context, req models.UserChangePasswordInputModel) (*handler.Response[uuid.UUID], error) {
	user, err := u.GetUser(c)
	if err != nil {
		return nil, err
	}

	err = u.userService.ChangePassword(c, user.ID, &req)

	if err != nil {
		return nil, err
	}

	return handler.OK(&user.ID, "Password changed successfully"), nil
}

func (u *ProfileHandler) ChangeProfile(c *gin.Context, req models.UserChangeProfileInputModel) (*handler.Response[models.UserProfileOutPutModel], error) {
	user, err := u.GetUser(c)
	if err != nil {
		return nil, err
	}

	req.ID = user.ID

	err = u.userService.ChangeProfile(c, user.ID, &req)

	if err != nil {
		return nil, err
	}

	return u.Profile(c)
}
