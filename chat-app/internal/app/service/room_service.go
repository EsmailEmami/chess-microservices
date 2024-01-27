package service

import (
	"context"

	appModels "github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/esmailemami/chess/chat/internal/models"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	sharedModels "github.com/esmailemami/chess/shared/models"
	"github.com/esmailemami/chess/shared/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomService struct {
	service.BaseService[models.Room]

	userService *service.UserService
}

func NewRoomService(userService *service.UserService) *RoomService {
	return &RoomService{
		userService: userService,
	}
}

func (r *RoomService) CreatePrivateRoom(ctx context.Context, currentUser *sharedModels.User, req *appModels.CreatePrivateRoomInputModel) (*models.Room, error) {

	if currentUser.ID == req.UserID {
		return nil, errs.BadRequestErr().Msg("The requested user is not valied")
	}

	db := psql.DBContext(ctx)
	tx := db.Begin()

	room := req.ToDBModel()

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

func (r *RoomService) CreateGlobalRoomRoom(ctx context.Context, currentUser *sharedModels.User, req *appModels.CreateGlobalRoomInputModel) (*models.Room, error) {
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

func (r *RoomService) OpenRoom(ctx context.Context, id uuid.UUID) (*appModels.RoomOutPutModel, error) {
	db := psql.DBContext(ctx)

	roomDB := db.Model(&models.Room{}).
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Model(&models.UserRoom{})
		}).
		Preload("Users.User", func(db *gorm.DB) *gorm.DB {
			return db.Model(&sharedModels.User{})
		}).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Model(&models.Message{})
		}).
		Preload("Messages.ReplyTo", func(db *gorm.DB) *gorm.DB {
			return db.Model(&models.Message{})
		})

	var room appModels.RoomOutPutModel

	if err := roomDB.First(&room, "id=?", id).Error; err != nil {
		return nil, errs.NotFoundErr().WithError(err)
	}

	return &room, nil
}
