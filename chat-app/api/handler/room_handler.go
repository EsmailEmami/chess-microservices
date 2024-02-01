package handler

import (
	"github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/app/service"
	chatroom "github.com/esmailemami/chess/chat/internal/chat-room"
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

// GetRooms godoc
// @Tags room
// @Accept json
// @Produce json
// @Security Bearer
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} handler.Response[handler.ListResponse[models.RoomsOutPutModel]]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} map[string]any
// @Router /room [get]
func (r *RoomHandler) GetRooms(ctx *gin.Context, params models.RoomQueryParams) (*handler.Response[handler.ListResponse[models.RoomsOutPutModel]], error) {
	rooms, totalRecords, err := r.roomService.GetRooms(ctx, &params)

	if err != nil {
		return nil, err
	}

	return handler.ListOK[models.RoomsOutPutModel](params.Page, params.Limit, totalRecords, rooms), nil
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

	// join to the web socket rooms if user is online
	chatroom.ConnectPrvateRoom(room.ID, currentUser.ID)
	chatroom.ConnectPrvateRoom(room.ID, req.UserID)

	return handler.OK(&room.ID), nil
}

func (r *RoomHandler) CreatePublicRoom(ctx *gin.Context, req models.CreatePublicRoomInputModel) (*handler.Response[uuid.UUID], error) {
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	currentUser, err := r.GetUser(ctx)
	if err != nil {
		return nil, errs.BadRequestErr()
	}

	room, err := r.roomService.CreatePublicRoomRoom(ctx, currentUser, &req)

	if err != nil {
		return nil, err
	}

	// join to the web socket rooms if user is online
	chatroom.ConnectPublicRoom(room.ID, currentUser.ID)
	for _, userID := range req.Users {
		chatroom.ConnectPublicRoom(room.ID, userID)
	}

	return handler.OK(&room.ID), nil
}

func (r *RoomHandler) JoinRoom(ctx *gin.Context, id uuid.UUID) (*handler.Response[bool], error) {
	currentUser, err := r.GetUser(ctx)
	if err != nil {
		return nil, errs.BadRequestErr()
	}

	err = r.roomService.JoinRoom(ctx, id, currentUser.ID)

	if err != nil {
		return nil, err
	}

	// join to the web socket rooms if user is online
	chatroom.JoinRoom(id, currentUser)

	return handler.OKBool(), nil
}

func (r *RoomHandler) LeftRoom(ctx *gin.Context, id uuid.UUID) (*handler.Response[bool], error) {
	currentUser, err := r.GetUser(ctx)
	if err != nil {
		return nil, errs.BadRequestErr()
	}

	err = r.roomService.LeftRoom(ctx, id, currentUser.ID)

	if err != nil {
		return nil, err
	}

	// join to the web socket rooms if user is online
	chatroom.LeftRoom(id, currentUser)

	return handler.OKBool(), nil
}
