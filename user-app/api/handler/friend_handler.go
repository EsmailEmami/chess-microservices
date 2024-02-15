package handler

import (
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/handler"
	"github.com/esmailemami/chess/user/internal/app/models"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FriendHandler struct {
	handler.Handler

	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
	}
}

// MakeFriend godoc
// @Tags friend
// @Accept json
// @Produce json
// @Security Bearer
// @Param friendId   path  string  true  "friend id"
// @Success 200 {object} handler.Response[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /friend/make/{friendId} [post]
func (h *FriendHandler) MakeFriend(ctx *gin.Context, friendID uuid.UUID) (*handler.Response[uuid.UUID], error) {
	currentUser := h.GetUser(ctx)
	if currentUser == nil {
		return nil, errs.UnAuthorizedErr()
	}

	friend, err := h.friendService.MakeFriend(ctx, currentUser.ID, friendID)

	if err != nil {
		return nil, err
	}

	return handler.OK(&friend.ID), nil
}

// RemoveFriend godoc
// @Tags friend
// @Accept json
// @Produce json
// @Security Bearer
// @Param friendId   path  string  true  "friend id"
// @Success 200 {object} handler.Response[bool]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /friend/remove/{friendId} [post]
func (h *FriendHandler) RemoveFriend(ctx *gin.Context, friendID uuid.UUID) (*handler.Response[bool], error) {
	currentUser := h.GetUser(ctx)
	if currentUser == nil {
		return nil, errs.UnAuthorizedErr()
	}

	err := h.friendService.RemoveFriend(ctx, currentUser.ID, friendID)

	if err != nil {
		return nil, err
	}

	return handler.OKBool(), nil
}

// GetFriends godoc
// @Tags friend
// @Accept json
// @Produce json
// @Security Bearer
// @Param searchTern   query  string  false  "search term"
// @Success 200 {object} handler.Response[[]models.FriendOutPutModel]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /friend [get]
func (h *FriendHandler) GetFriends(ctx *gin.Context, params models.FriendQueryParams) (*handler.Response[[]models.FriendOutPutModel], error) {
	currentUser := h.GetUser(ctx)
	if currentUser == nil {
		return nil, errs.UnAuthorizedErr()
	}

	friends, err := h.friendService.GetFriends(ctx, currentUser.ID, &params)

	if err != nil {
		return nil, err
	}

	return handler.OK(&friends), nil
}
