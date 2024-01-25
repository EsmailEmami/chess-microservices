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
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) Login(c *gin.Context, req models.LoginInputModel) (*handler.Response[models.LoginOutputModel], error) {
	// validate the model
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	resp, err := a.authService.Login(c, &req)

	if err != nil {
		return nil, err
	}

	return handler.OK(resp), nil
}

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
