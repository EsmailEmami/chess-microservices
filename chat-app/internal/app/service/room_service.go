package service

import (
	"context"
	"time"

	appModels "github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/models"
	"github.com/esmailemami/chess/chat/internal/util"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/logging"
	sharedModels "github.com/esmailemami/chess/shared/models"
	"github.com/esmailemami/chess/shared/service"
	"github.com/esmailemami/chess/shared/util/dbutil"
	"github.com/google/uuid"
)

const (
	roomCacheDuration = 12 * time.Hour
)

type RoomService struct {
	service.BaseService[models.Room]

	cache *redis.Redis
}

func NewRoomService(cache *redis.Redis) *RoomService {
	return &RoomService{
		cache: cache,
	}
}

func (r *RoomService) CreatePrivateRoom(ctx context.Context, currentUser *sharedModels.User, req *appModels.CreatePrivateRoomInputModel) (*models.Room, error) {

	if currentUser.ID == req.UserID {
		return nil, errs.BadRequestErr().Msg("The requested user is not valied")
	}

	db := psql.DBContext(ctx)
	tx := db.Begin()

	var (
		room1 = currentUser.ID.String() + "_" + req.UserID.String()
		room2 = req.UserID.String() + "_" + currentUser.ID.String()
	)

	if dbutil.Exists(&models.Room{}, "(name = ? || name = ?) AND is_private = ?", room1, room2, true) {
		tx.Rollback()
		return nil, errs.BadRequestErr().Msg("room already exists")
	}

	room := req.ToDBModel()
	room.Name = room1

	if err := tx.Create(room).Error; err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	// user rooms
	userRooms := []models.UserRoom{
		{
			UserID: currentUser.ID,
			RoomID: room.ID,
			BaseModel: sharedModels.BaseModel{
				ID: uuid.New(),
			},
		},
		{
			UserID: req.UserID,
			RoomID: room.ID,
			BaseModel: sharedModels.BaseModel{
				ID: uuid.New(),
			},
		},
	}

	if err := tx.Model(&models.UserRoom{}).Create(&userRooms).Error; err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return room, nil
}

func (r *RoomService) CreatePublicRoomRoom(ctx context.Context, currentUser *sharedModels.User, req *appModels.CreatePublicRoomInputModel) (*models.Room, error) {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	room := req.ToDBModel()

	if err := tx.Create(room).Error; err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	// user rooms
	userRooms := make([]models.UserRoom, len(req.Users)+1)
	userRooms[0] = models.UserRoom{
		UserID: currentUser.ID,
		RoomID: room.ID,
		BaseModel: sharedModels.BaseModel{
			ID: uuid.New(),
		},
	}

	for i, userId := range req.Users {
		userRooms[i+1] = models.UserRoom{
			UserID: userId,
			RoomID: room.ID,
			BaseModel: sharedModels.BaseModel{
				ID: uuid.New(),
			},
		}
	}

	if err := tx.Model(&models.UserRoom{}).Create(&userRooms).Error; err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return room, nil
}

func (r *RoomService) GetUserRoomIDs(ctx context.Context, userID uuid.UUID, loadPrivate bool) ([]uuid.UUID, error) {
	db := psql.DBContext(ctx)

	var roomIDs []uuid.UUID

	if err := db.Model(&models.UserRoom{}).
		Joins("INNER JOIN chat.room on chat.room.id = chat.user_room.room_id").
		Where("chat.user_room.user_id = ? AND chat.room.is_private = ? AND chat.room.deleted_at IS NULL", userID, loadPrivate).Select("chat.user_room.room_id").Find(&roomIDs).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return roomIDs, nil
}

func (r *RoomService) Get(ctx context.Context, id uuid.UUID, userID *uuid.UUID) (*appModels.RoomOutPutModel, error) {
	room := &appModels.RoomOutPutModel{}

	if err := r.cache.UnmarshalToObject(r.getRoomCacheKey(id), room); err == nil {
		return room, nil
	}

	room, err := r.setRoomCache(ctx, id)
	if err != nil {
		return room, err
	}

	// check for profile
	if !room.IsPrivate || userID == nil {
		return room, nil
	}

	avatar := ""
	for _, user := range room.Users {
		if user.ID != *userID {
			avatar = user.Profile
			break
		}
	}
	room.Avatar = avatar

	return room, nil
}

func (r *RoomService) EditRoom(ctx context.Context, id uuid.UUID, req *appModels.EditRoomModel) error {
	db := psql.DBContext(ctx)

	var room models.Room

	if err := db.First(&room, "id = ?", id).Error; err != nil {
		return errs.NotFoundErr().WithError(err)
	}

	if room.IsPrivate {
		return errs.BadRequestErr().Msg("room is not public")
	}

	req.MergeWithDbModel(&room)

	if err := db.Save(&room).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	// reset the cache
	if _, err := r.setRoomCache(ctx, id); err != nil {
		logging.ErrorE("failed to reset room cache", err)
	}

	return nil
}

func (r *RoomService) JoinRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	db := psql.DBContext(ctx)

	var room models.Room

	if err := db.Model(&models.Room{}).First(&room, "id = ?", roomID).Error; err != nil {
		return errs.NotFoundErr().WithError(err).Msg("room not found")
	}

	if room.IsPrivate {
		return errs.BadRequestErr().Msg("room is not public")
	}

	if dbutil.Exists(&models.UserRoom{}, "user_id = ? AND room_id = ?", userID, roomID) {
		return errs.BadRequestErr().Msg("user already joined")
	}

	dbModel := models.UserRoom{
		UserID: userID,
		RoomID: roomID,
	}

	if err := db.Create(&dbModel).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	// reset the cache
	if _, err := r.setRoomCache(ctx, roomID); err != nil {
		logging.ErrorE("failed to reset room cache", err)
	}

	return nil
}

func (r *RoomService) GetRooms(ctx context.Context, params *appModels.RoomQueryParams) (result []appModels.RoomsOutPutModel, totalRecords int64, err error) {
	db := psql.DBContext(ctx)
	qry := db.Model(&models.Room{}).Where("is_private = ?", false)

	qry = dbutil.Filter(qry, params.SearchTerm, "name")
	totalRecords, err = dbutil.Paginate(qry, params.Page, params.Limit, &result)
	return
}

func (r *RoomService) LeftRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	db := psql.DBContext(ctx)

	var room models.Room

	if err := db.Model(&models.Room{}).First(&room, "id = ?", roomID).Error; err != nil {
		return errs.NotFoundErr().WithError(err).Msg("room not found")
	}

	if room.IsPrivate {
		return errs.BadRequestErr().Msg("room is not public")
	}

	if !dbutil.Exists(&models.UserRoom{}, "user_id = ? AND room_id = ?", userID, roomID) {
		return errs.BadRequestErr().Msg("user is not joined")
	}

	if err := db.Model(&models.UserRoom{}).Where("user_id = ? AND room_id = ?", userID, roomID).Delete(&models.UserRoom{}).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	// reset the cache
	if _, err := r.setRoomCache(ctx, roomID); err != nil {
		logging.ErrorE("failed to reset room cache", err)
	}

	return nil
}

func (r *RoomService) Delete(ctx context.Context, currentUser *sharedModels.User, id uuid.UUID) error {
	db := psql.DBContext(ctx)

	room, err := r.getRoom(ctx, id)
	if err != nil {
		return err
	}

	if !r.hasPermission(room, currentUser) && !room.IsPrivate {
		return errs.AccessDeniedError().Msg("you do not have permission of deleting this room")
	}

	if err := db.Where("id = ?", id).Delete(&models.Room{}).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	// delete the cache
	_ = r.cache.Delete(r.getRoomCacheKey(id))

	return nil
}

func (r *RoomService) UpdateAvatar(ctx context.Context, roomID uuid.UUID, avatar string) error {
	db := psql.DBContext(ctx)

	room, err := r.getRoom(ctx, roomID)
	if err != nil {
		return err
	}

	if room.IsPrivate {
		return errs.BadRequestErr().Msg("room is not public")
	}

	if err := db.Model(&models.Room{}).Where("id = ?", roomID).UpdateColumn("avatar", avatar).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	// delete the cache
	_ = r.cache.Delete(r.getRoomCacheKey(roomID))

	return nil
}

func (r *RoomService) DeleteCache(id uuid.UUID) {
	err := r.cache.Delete(r.getRoomCacheKey(id))

	if err != nil {
		logging.ErrorE("failed to remove room cache", err)
	}
}

func (r *RoomService) hasPermission(room *appModels.RoomOutPutModel, user *sharedModels.User) bool {
	if user.IsAdmin() {
		return true
	}

	// room is public and just need to check the creator
	if !room.IsPrivate && (room.CreatedByID != nil && *room.CreatedByID == user.ID) {
		return true
	}

	if !room.IsPrivate {
		return false
	}

	hasAccess := false

	// check the private room users
	for _, roomUser := range room.Users {
		if roomUser.ID == user.ID {
			hasAccess = true
			break
		}
	}

	return hasAccess
}

func (r *RoomService) getRoom(ctx context.Context, id uuid.UUID) (*appModels.RoomOutPutModel, error) {
	var room appModels.RoomOutPutModel

	if err := r.cache.UnmarshalToObject(r.getRoomCacheKey(id), &room); err == nil {
		newRoom, err := r.setRoomCache(ctx, id)

		if err != nil {
			return nil, errs.NotFoundErr()
		}

		return newRoom, nil
	}

	return &room, nil
}

func (r *RoomService) getRoomCacheKey(id uuid.UUID) string {
	return "room_" + id.String()
}

func (r *RoomService) setRoomCache(ctx context.Context, id uuid.UUID) (*appModels.RoomOutPutModel, error) {
	db := psql.DBContext(ctx)

	var dbRoom models.Room

	if err := db.Model(&models.Room{}).
		Preload("Users").
		Preload("Users.User").First(&dbRoom, "id = ?", id).Error; err != nil {
		return nil, errs.NotFoundErr().WithError(err)
	}

	room := appModels.RoomOutPutModel{
		ID:          dbRoom.ID,
		Name:        dbRoom.Name,
		IsPrivate:   dbRoom.IsPrivate,
		CreatedByID: dbRoom.CreatedByID,
		Avatar:      util.FilePathPrefix(dbRoom.Avatar),
		Users:       make([]appModels.RoomUserOutPutModel, len(dbRoom.Users)),
	}

	for i, userRoom := range dbRoom.Users {
		room.Users[i] = appModels.RoomUserOutPutModel{
			ID:        userRoom.UserID,
			FirstName: userRoom.User.FirstName,
			LastName:  userRoom.User.LastName,
			Username:  userRoom.User.Username,
			Profile:   util.FilePathPrefix(userRoom.User.Profile),
		}
	}

	if err := r.cache.Set(r.getRoomCacheKey(id), &room, roomCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return &room, nil
}
