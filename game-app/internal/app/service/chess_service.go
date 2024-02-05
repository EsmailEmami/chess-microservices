package service

import (
	"context"
	"time"

	appModels "github.com/esmailemami/chess/game/internal/app/models"
	"github.com/esmailemami/chess/game/internal/models"
	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/errs"
	sharedModels "github.com/esmailemami/chess/shared/models"
	"github.com/esmailemami/chess/shared/service"
	"github.com/esmailemami/chess/shared/util"
	"github.com/google/uuid"
)

var (
	chessCacheDuration        = 20 * time.Minute
	chessWatcherCacheDuration = 5 * time.Hour
)

type ChessService struct {
	service.BaseService[models.Chess]

	userService *service.UserService
	cache       *redis.Redis
}

func NewChessService(cache *redis.Redis, userService *service.UserService) *ChessService {
	return &ChessService{
		cache:       cache,
		userService: userService,
	}
}

func (*ChessService) GetGamesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Chess, error) {
	var Chesss []models.Chess

	db := psql.DBContext(ctx)

	if err := db.Model(&models.Chess{}).Where("white_player_id=? OR black_player_id=?", userID, userID).
		Find(&Chesss).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return Chesss, nil
}

func (*ChessService) GetChessIDsByUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var Chesss []uuid.UUID

	db := psql.DBContext(ctx)

	if err := db.Model(&models.Chess{}).Where("white_player_id=? OR black_player_id=?", userID, userID).
		Select("id").
		Find(&Chesss).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return Chesss, nil
}

func (g *ChessService) Get(ctx context.Context, id uuid.UUID) (*appModels.ChessOutputModel, error) {
	var output appModels.ChessOutputModel

	if err := g.cache.UnmarshalToObject(g.getChessCacheKey(id), &output); err == nil {
		return &output, nil
	}

	db := psql.DBContext(ctx)

	var model models.Chess

	if err := db.Preload("WhitePlayer").Preload("BlackPlayer").Find(&model, "id=?", id).Error; err != nil {
		return nil, errs.NotFoundErr()
	}

	return g.setChessCache(ctx, id)
}

func (g *ChessService) JoinGame(ctx context.Context, currentUser *sharedModels.User, id uuid.UUID) error {
	db := psql.DBContext(ctx)

	var chess models.Chess

	if err := db.Find(&chess, "id=?", id).Error; err != nil {
		return errs.NotFoundErr()
	}

	if chess.Status != models.ChessStatusWaiting {
		return errs.BadRequestErr().Msg("you can not join the Chess")
	}

	if chess.WhitePlayerID == nil {
		chess.WhitePlayerID = &currentUser.ID
	} else {
		chess.BlackPlayerID = &currentUser.ID
	}

	chess.Status = models.ChessStatusOpen

	if err := db.Save(&chess).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (g *ChessService) MoveChessPiece(ctx context.Context, id uuid.UUID, piece *chessboard.Piece, from, to chessboard.Position) error {
	db := psql.DBContext(ctx)

	var chess models.Chess

	if err := db.First(&chess, "id = ?", id).Error; err != nil {
		return errs.NotFoundErr().WithError(err)
	}

	for i, piece := range chess.Pieces {
		if !(piece.Row == from.Row && piece.Col == from.Col) && !(piece.Row == to.Row && piece.Col == to.Col) {
			continue
		}

		chess.Pieces = util.ArrayRemoveIndex[models.ChessPiece](chess.Pieces, i)
	}

	player := models.GetChessPlayerFromColor(piece.Color)

	chess.Pieces = append(chess.Pieces, models.ChessPiece{
		Piece:  string(piece.Type),
		Row:    to.Row,
		Col:    to.Col,
		Player: player,
	})

	chess.Moves = append(chess.Moves, models.ChessMove{
		Player: player,
		From:   from.String(),
		To:     to.String(),
	})

	chess.SwitchTurn()

	if err := db.Save(&chess).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (g *ChessService) Chectmate(ctx context.Context, id uuid.UUID, winnerID uuid.UUID) error {
	db := psql.DBContext(ctx)

	var chess models.Chess

	if err := db.First(&chess, "id = ?", id).Error; err != nil {
		return errs.NotFoundErr().WithError(err)
	}

	chess.Status = models.ChessStatusClose
	chess.WinnerID = &winnerID

	if err := db.Save(&chess).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (g *ChessService) setChessCache(ctx context.Context, id uuid.UUID) (*appModels.ChessOutputModel, error) {
	db := psql.DBContext(ctx)

	var chess models.Chess

	if err := db.Preload("WhitePlayer").Preload("BlackPlayer").Find(&chess, "id=?", id).Error; err != nil {
		return nil, errs.NotFoundErr()
	}

	output := &appModels.ChessOutputModel{
		ID:            chess.ID,
		Turn:          chess.Turn,
		Moves:         chess.Moves,
		Pieces:        chess.Pieces,
		Status:        chess.Status,
		Winner:        chess.WinnerID,
		WhitePlayerID: chess.WhitePlayerID,
		BlackPlayerID: chess.BlackPlayerID,
	}

	if chess.WhitePlayer != nil {
		output.WhitePlayer = &appModels.ChessPlayerOutputModel{
			ID:        *chess.WhitePlayerID,
			FirstName: chess.WhitePlayer.FirstName,
			LastName:  chess.WhitePlayer.LastName,
			Username:  chess.WhitePlayer.Username,
		}
	}

	if chess.BlackPlayer != nil {
		output.BlackPlayer = &appModels.ChessPlayerOutputModel{
			ID:        *chess.BlackPlayerID,
			FirstName: chess.BlackPlayer.FirstName,
			LastName:  chess.BlackPlayer.LastName,
			Username:  chess.BlackPlayer.Username,
		}
	}

	// cache the data
	if err := g.cache.Set(g.getChessCacheKey(id), output, chessCacheDuration); err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return output, nil
}

func (g *ChessService) GetUserWatchingGames(userID uuid.UUID) []uuid.UUID {
	var chessIDs []uuid.UUID
	_ = g.cache.UnmarshalToObject(g.getUserChessWatcherCacheKey(userID), &chessIDs)
	return chessIDs
}

func (g *ChessService) SetWatcher(chessID, userID uuid.UUID) error {
	var chessIDs []uuid.UUID

	_ = g.cache.UnmarshalToObject(g.getUserChessWatcherCacheKey(userID), &chessIDs)

	chessIDs = append(chessIDs, chessID)

	return g.cache.Set(g.getUserChessWatcherCacheKey(userID), &chessIDs, chessWatcherCacheDuration)
}

func (g *ChessService) NewChess(ctx context.Context, currentUser *sharedModels.User, req *appModels.CreateChessInputModel) (*models.Chess, error) {
	db := psql.DBContext(ctx)

	chessboard := chessboard.NewDefault()

	var (
		whitePlayer, blackPlayer *sharedModels.User
	)

	if req.Color == "white" {
		whitePlayer = currentUser
	} else {
		blackPlayer = currentUser
	}

	if req.PlayingWith != nil {

		opponetUser, err := g.userService.Get(ctx, *req.PlayingWith)

		if err != nil {
			return nil, err
		}

		if req.Color == "white" {
			blackPlayer = opponetUser
		} else {
			whitePlayer = opponetUser
		}
	}

	chess := models.NewChess(whitePlayer, blackPlayer, chessboard.GetPieces())

	if err := db.Create(chess).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return chess, nil
}

func (g *ChessService) DeleteWatcherCache(userID uuid.UUID) error {
	return g.cache.Delete(g.getUserChessWatcherCacheKey(userID))
}

func (g *ChessService) getUserChessWatcherCacheKey(userID uuid.UUID) string {
	return "chess_watcher_" + userID.String()
}

func (g *ChessService) getChessCacheKey(id uuid.UUID) string {
	return "chess_" + id.String()
}
