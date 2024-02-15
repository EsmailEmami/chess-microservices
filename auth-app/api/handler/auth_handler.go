package handler

import (
	"github.com/esmailemami/chess/auth/internal/app/models"
	"github.com/esmailemami/chess/auth/internal/app/service"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	handler.Handler

	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param input   body  models.LoginInputModel  true  "input model"
// @Success 200 {object} handler.Response[models.LoginOutputModel]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /login [post]
func (a *AuthHandler) Login(c *gin.Context, req models.LoginInputModel) (*handler.Response[models.LoginOutputModel], error) {
	// validate the model
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	resp, err := a.authService.Login(c, &req)

	if err != nil {
		return nil, err
	}

	c.SetCookie("Authorization", resp.Token, int(resp.ExpiresIn), "/", "", true, true)

	return handler.OK(resp), nil
}

// Register godoc
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param input   body  models.RegisterInputModel  true  "input model"
// @Success 200 {object} handler.Response[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /register [post]
func (a *AuthHandler) Register(c *gin.Context, req models.RegisterInputModel) (*handler.Response[uuid.UUID], error) {
	// validate the model
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	resp, err := a.authService.Register(c, &req)

	if err != nil {
		return nil, err
	}

	return handler.OK(&resp.ID), nil
}

// Logout godoc
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} handler.Response[bool]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /logout [post]
func (a *AuthHandler) Logout(c *gin.Context) (*handler.Response[bool], error) {
	jwtID := c.GetHeader("JwtId")

	if jwtID != "" {
		a.authService.RevokeToken(c, uuid.MustParse(jwtID))
	}

	a.removeAuthorizationCookie(c)

	return handler.OKBool(), nil
}

func (a *AuthHandler) removeAuthorizationCookie(ctx *gin.Context) {
	ctx.SetCookie("Authorization", "", 0, "/", "", true, true)
}
