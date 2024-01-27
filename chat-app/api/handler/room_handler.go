package handler

import (
	"github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoomHandler struct {
	handler.Handler

	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func (r *RoomHandler) CreatePrivateRoom(ctx *gin.Context, req models.CreatePrivateRoomInputModel) (*handler.Response[uuid.UUID], error) {
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	currentUser, err := r.GetUser(ctx)
	if err != nil {
		return nil, errs.BadRequestErr()
	}

	room, err := r.roomService.CreatePrivateRoom(ctx, currentUser, &req)

	if err != nil {
		return nil, err
	}

	return handler.OK(&room.ID), nil
}

func (r *RoomHandler) CreateGlobalRoom(ctx *gin.Context, req models.CreateGlobalRoomInputModel) (*handler.Response[uuid.UUID], error) {
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	currentUser, err := r.GetUser(ctx)
	if err != nil {
		return nil, errs.BadRequestErr()
	}

	room, err := r.roomService.CreateGlobalRoomRoom(ctx, currentUser, &req)

	if err != nil {
		return nil, err
	}

	return handler.OK(&room.ID), nil
}

func (r *RoomHandler) GetRoom(ctx *gin.Context, id uuid.UUID) (*handler.Response[models.RoomOutPutModel], error) {
	room, err := r.roomService.OpenRoom(ctx, id)

	if err != nil {
		return nil, err
	}

	return handler.OK(room), nil
}
