package handler

import (
	"github.com/esmailemami/chess/shared/errs"
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

// Profile godoc
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} handler.Response[models.UserProfileOutPutModel]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /profile [get]
func (u *ProfileHandler) Profile(c *gin.Context) (*handler.Response[models.UserProfileOutPutModel], error) {
	user := u.GetUser(c)
	if user == nil {
		return nil, errs.UnAuthorizedErr()
	}

	profile, err := u.userService.GetProfile(c, user.ID)

	if err != nil {
		return nil, err
	}

	return handler.OK(profile), nil
}

// ChangePassword godoc
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param input   body  models.UserChangePasswordInputModel  true  "input model"
// @Success 200 {object} handler.Response[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /profile/change-password [post]
func (u *ProfileHandler) ChangePassword(c *gin.Context, req models.UserChangePasswordInputModel) (*handler.Response[uuid.UUID], error) {
	user := u.GetUser(c)
	if user == nil {
		return nil, errs.UnAuthorizedErr()
	}

	err := u.userService.ChangePassword(c, user.ID, &req)

	if err != nil {
		return nil, err
	}

	return handler.OK(&user.ID, "Password changed successfully"), nil
}

// ChangeProfile godoc
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param input   body  models.UserChangeProfileInputModel  true  "input model"
// @Success 200 {object} handler.Response[models.UserProfileOutPutModel]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /profile [post]
func (u *ProfileHandler) ChangeProfile(c *gin.Context, req models.UserChangeProfileInputModel) (*handler.Response[models.UserProfileOutPutModel], error) {
	user := u.GetUser(c)
	if user == nil {
		return nil, errs.UnAuthorizedErr()
	}
	req.ID = user.ID

	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	err := u.userService.ChangeProfile(c, user.ID, &req)

	if err != nil {
		return nil, err
	}

	return u.Profile(c)
}
