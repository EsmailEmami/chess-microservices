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
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Param page  query  string  false  "page size"
// @Param limit  query  string  false  "length of records to show"
// @Param searchTerm  query  string  false  "search for item"
// @Success 200 {object} handler.JSONResponse[handler.ListResponse[models.RoomsOutPutModel]]
// @Router /room [get]
func (r *RoomHandler) GetRooms(ctx *gin.Context, params models.RoomQueryParams) (handler.Response, error) {
	rooms, totalRecords, err := r.roomService.GetRooms(ctx, &params)

	if err != nil {
		return nil, err
	}

	return handler.ListOK(params.Page, params.Limit, totalRecords, rooms), nil
}

// CreatePrivateRoom godoc
// @Tags room
// @Accept json
// @Produce json
// @Security Bearer
// @Param input   body  models.CreatePrivateRoomInputModel  true  "input model"
// @Success 200 {object} handler.JSONResponse[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /room/private [post]
func (r *RoomHandler) CreatePrivateRoom(ctx *gin.Context, req models.CreatePrivateRoomInputModel) (handler.Response, error) {
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	currentUser := r.GetUser(ctx)
	if currentUser == nil {
		return nil, errs.UnAuthorizedErr()
	}

	room, err := r.roomService.CreatePrivateRoom(ctx, currentUser, &req)

	if err != nil {
		return nil, err
	}

	// join to the web socket rooms if user is online
	chatroom.ConnectPrivateRoom(room.ID, currentUser.ID)
	chatroom.ConnectPrivateRoom(room.ID, req.UserID)

	return handler.OK(&room.ID), nil
}

// CreatePublicRoom godoc
// @Tags room
// @Accept json
// @Produce json
// @Security Bearer
// @Param input   body  models.CreatePublicRoomInputModel  true  "input model"
// @Success 200 {object} handler.JSONResponse[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /room/public [post]
func (r *RoomHandler) CreatePublicRoom(ctx *gin.Context, req models.CreatePublicRoomInputModel) (handler.Response, error) {
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	currentUser := r.GetUser(ctx)
	if currentUser == nil {
		return nil, errs.UnAuthorizedErr()
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

// JoinRoom godoc
// @Tags room
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Success 200 {object} handler.JSONResponse[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /room/join/{id} [post]
func (r *RoomHandler) JoinRoom(ctx *gin.Context, id uuid.UUID) (handler.Response, error) {
	currentUser := r.GetUser(ctx)
	if currentUser == nil {
		return nil, errs.UnAuthorizedErr()
	}

	err := r.roomService.JoinRoom(ctx, id, currentUser.ID)

	if err != nil {
		return nil, err
	}

	// join to the web socket rooms if user is online
	chatroom.JoinRoom(id, currentUser)

	return handler.OKBool(), nil
}

// LeftRoom godoc
// @Tags room
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Success 200 {object} handler.JSONResponse[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /room/left/{id} [post]
func (r *RoomHandler) LeftRoom(ctx *gin.Context, id uuid.UUID) (handler.Response, error) {
	currentUser := r.GetUser(ctx)
	if currentUser == nil {
		return nil, errs.UnAuthorizedErr()
	}

	err := r.roomService.LeftRoom(ctx, id, currentUser.ID)

	if err != nil {
		return nil, err
	}

	// join to the web socket rooms if user is online
	chatroom.LeftRoom(id, currentUser)

	return handler.OKBool(), nil
}

// DeleteRoom godoc
// @Tags room
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Success 200 {object} handler.JSONResponse[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /room/delete/{id} [post]
func (r *RoomHandler) DeleteRoom(ctx *gin.Context, id uuid.UUID) (handler.Response, error) {
	user := r.GetUser(ctx)
	if user == nil {
		return nil, errs.UnAuthorizedErr()
	}

	err := r.roomService.Delete(ctx, user, id)

	if err != nil {
		return nil, err
	}

	// delete room from socket
	chatroom.DeleteRoom(id)

	return handler.OKBool(), nil
}

// EditRoom godoc
// @Tags room
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Param input   body  models.EditRoomModel  true  "input model"
// @Success 200 {object} handler.JSONResponse[bool]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /room/edit/{id} [post]
func (r *RoomHandler) EditRoom(ctx *gin.Context, id uuid.UUID, req models.EditRoomModel) (handler.Response, error) {
	if err := req.Validate(); err != nil {
		return nil, errs.ValidationErr(err)
	}

	if err := r.roomService.EditRoom(ctx, id, &req); err != nil {
		return nil, err
	}

	// jnotify the connected clients to get last update of room
	chatroom.RoomEdited(id)

	return handler.OKBool(), nil
}
