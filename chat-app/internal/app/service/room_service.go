package service

import (
	"context"
	"time"

	appModels "github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/models"
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

func (r *RoomService) Get(ctx context.Context, id uuid.UUID) (*appModels.RoomOutPutModel, error) {
	var room appModels.RoomOutPutModel

	if err := r.cache.UnmarshalToObject(r.getRoomCacheKey(id), &room); err == nil {
		return &room, nil
	}

	return r.setRoomCache(ctx, id)
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
		ID:        dbRoom.ID,
		Name:      dbRoom.Name,
		IsPrivate: dbRoom.IsPrivate,
		Users:     make([]appModels.RoomUserOutPutModel, len(dbRoom.Users)),
	}

	for i, userRoom := range dbRoom.Users {
		room.Users[i] = appModels.RoomUserOutPutModel{
			ID:        userRoom.UserID,
			FirstName: userRoom.User.FirstName,
			LastName:  userRoom.User.LastName,
			Username:  userRoom.User.Username,
		}
	}

	if err := r.cache.Set(r.getRoomCacheKey(id), &room, roomCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return &room, nil
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

func (r *RoomService) getRoomCacheKey(id uuid.UUID) string {
	return "room_" + id.String()
}

func (b *RoomService) Delete(ctx context.Context, id uuid.UUID) error {
	db := psql.DBContext(ctx)

	var room models.Room

	if err := db.Model(&models.Room{}).First(&room, "id = ?", id).Error; err != nil {
		return errs.NotFoundErr().WithError(err)
	}

	//TODO: check who is creator or admin!
	if err := db.Delete(&room).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (r *RoomService) UpdateAvatar(ctx context.Context, roomID uuid.UUID, avatar string) error {
	db := psql.DBContext(ctx)

	var room models.Room

	if err := db.Model(&models.Room{}).First(&room, "id = ?", roomID).Error; err != nil {
		return errs.NotFoundErr().WithError(err)
	}

	if room.IsPrivate {
		return errs.BadRequestErr().Msg("room is not public")
	}

	if err := db.Model(&models.Room{}).Where("id = ?", roomID).UpdateColumn("avatar", avatar).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}
