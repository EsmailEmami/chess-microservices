package service

import (
	"context"

	"github.com/esmailemami/chess/game/internal/models"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/service"
	"github.com/google/uuid"
)

type GameService struct {
	service.BaseService[models.Game]
}

func NewGameService() *GameService {
	return new(GameService)
}

func (*GameService) GetGamesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Game, error) {
	var games []models.Game

	db := psql.DBContext(ctx)

	if err := db.Model(&models.Game{}).Where("white_player_id=? OR black_player_id=?", userID, userID).
		Find(&games).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return games, nil
}

func (*GameService) GetGameIdsByUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var games []uuid.UUID

	db := psql.DBContext(ctx)

	if err := db.Model(&models.Game{}).Where("white_player_id=? OR black_player_id=?", userID, userID).
		Select("id").
		Find(&games).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return games, nil
}

func (b *GameService) Get(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	db := psql.DBContext(ctx)

	var model models.Game

	if err := db.Preload("WhitePlayer").Preload("BlackPlayer").Find(&model, "id=?", id).Error; err != nil {
		return nil, errs.NotFoundErr()
	}

	return &model, nil
}
