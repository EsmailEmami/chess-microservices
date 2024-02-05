package chess

import (
	"context"

	"github.com/esmailemami/chess/game/internal/app/models"
	"github.com/esmailemami/chess/game/internal/app/service"
	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/logging"
	sharedService "github.com/esmailemami/chess/shared/service"
	"github.com/google/uuid"
)

var games = make(map[uuid.UUID]*Board, 0)

func getBoard(ctx context.Context, gameID uuid.UUID) (*Board, error) {
	if board, ok := games[gameID]; ok {
		return board, nil
	}

	chessService := service.NewChessService(redis.GetConnection(), sharedService.NewUserService())

	game, err := chessService.Get(ctx, gameID)

	if err != nil {
		return nil, err
	}

	return loadGame(game, chessService)
}

func loadGame(chess *models.ChessOutputModel, chessService *service.ChessService) (*Board, error) {
	var (
		pieces = make([]*chessboard.ChessboardPiece, len(chess.Pieces))
		moves  = make([]*chessboard.ChessBoardMove, len(chess.Moves))
	)

	for i, piece := range chess.Pieces {
		pieces[i] = piece.ToChessPiece()
	}

	for i, move := range chess.Moves {

		fromPos, err := chessboard.GetPosition(move.From)

		if err != nil {
			logging.ErrorE("failed to parse game move", err)
			continue
		}

		toPos, err := chessboard.GetPosition(move.To)

		if err != nil {
			logging.ErrorE("failed to parse game move", err)
			continue
		}

		moves[i] = &chessboard.ChessBoardMove{
			From: *fromPos,
			To:   *toPos,
		}
	}

	chessboard := chessboard.New(pieces, moves)

	board, err := newBoard(chess.ID, chess.WhitePlayerID, chess.BlackPlayerID, chess.Status, chessboard, chessService)
	if err != nil {
		return nil, err
	}

	games[chess.ID] = board

	return board, nil
}

func deleteChess(chessID uuid.UUID) {
	delete(games, chessID)
}
