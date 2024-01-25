package game

import (
	"context"

	"github.com/esmailemami/chess/game/internal/models"
	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/game/pkg/websocket"
	sharedmodels "github.com/esmailemami/chess/shared/models"
	ws "github.com/esmailemami/chess/shared/websocket"
	"github.com/google/uuid"
)

func getBoard(ctx context.Context, gameID string) (*Board, error) {
	if board, ok := games[gameID]; ok {
		return board, nil
	}

	game, err := gameService.Get(ctx, uuid.MustParse(gameID))

	if err != nil {
		return nil, err
	}

	return loadGame(game)
}

func newGame(ctx context.Context, whitePlayerID, blackPlayerID *uuid.UUID) (*Board, error) {
	chessboard := chessboard.NewBoard()

	var (
		whitePlayer, blackPlayer *sharedmodels.User
		err                      error
	)

	if whitePlayerID != nil {
		whitePlayer, err = userService.Get(ctx, *whitePlayerID)
		if err != nil {
			return nil, err
		}
	}

	if blackPlayerID != nil {
		blackPlayer, err = userService.Get(ctx, *blackPlayerID)
		if err != nil {
			return nil, err
		}
	}

	game := models.NewGame(whitePlayer, blackPlayer, chessboard.GetPiecesPositions())

	if err := gameService.Create(ctx, game); err != nil {
		return nil, err
	}
	board, err := newBoard(game, chessboard)
	if err != nil {
		return nil, err
	}

	games[board.Game.ID.String()] = board

	return board, nil
}

func loadGame(game *models.Game) (*Board, error) {
	chessboard := chessboard.NewBoard(game.GetChessPieces()...)

	board, err := newBoard(game, chessboard)
	if err != nil {
		return nil, err
	}

	games[board.Game.ID.String()] = board

	return board, nil
}

func Run() {
	for {
		select {
		case req := <-websocket.NewGameCh:
			board, err := newGameRequest(req)
			if err != nil {
				websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
				continue
			}

			if board.Game.BlackPlayerID != nil {
				websocket.ChessWss.SendMessageToUser(*board.Game.BlackPlayerID, websocket.NewGame, board.ToOutput())
			}

			if board.Game.WhitePlayerID != nil {
				websocket.ChessWss.SendMessageToUser(*board.Game.WhitePlayerID, websocket.NewGame, board.ToOutput())
			}

		case req := <-websocket.GameValidMovesCh:
			moves, err := getBoardValidMovesRequest(req)

			if err != nil {
				websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
				continue
			}

			websocket.ChessWss.SendMessageToUser(req.UserID, websocket.GameValidMoves, moves)

		case req := <-websocket.GameBoardsCh:
			boards, err := userBoardsRequset(req)

			if err != nil {
				websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
				continue
			}

			websocket.ChessWss.SendMessageToUser(req.UserID, websocket.GameBoards, boards)

		case req := <-websocket.GameMovePieceCh:
			board, err := placePieceRequest(req)

			if err != nil {
				websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
				continue
			}

			websocket.ChessWss.SendMessageToUser(*board.Game.BlackPlayerID, websocket.GameBoardChanged, board.ToOutput())
			websocket.ChessWss.SendMessageToUser(*board.Game.WhitePlayerID, websocket.GameBoardChanged, board.ToOutput())

		case req := <-websocket.JoinGameCh:
			board, err := joinGameRequest(req)

			if err != nil {
				websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
				continue
			}

			if *board.WhitePlayerUserID == req.UserID {
				websocket.ChessWss.SendMessageToUser(*board.Game.WhitePlayerID, websocket.NewGame, board.ToOutput())
				websocket.ChessWss.SendMessageToUser(*board.Game.BlackPlayerID, websocket.GameBoardChanged, board.ToOutput())
			} else {
				websocket.ChessWss.SendMessageToUser(*board.Game.BlackPlayerID, websocket.NewGame, board.ToOutput())
				websocket.ChessWss.SendMessageToUser(*board.Game.WhitePlayerID, websocket.GameBoardChanged, board.ToOutput())
			}
		}
	}
}

func getBoardValidMovesRequest(req *ws.ClientMessage[websocket.GameValidMovesRequest]) ([]chessboard.Position, error) {
	board, err := getBoard(req.Ctx, req.Data.GameID)

	if err != nil {
		return nil, err
	}

	position, err := chessboard.GetPosition(req.Data.Position)

	if err != nil {
		return nil, err
	}

	return board.GetValidMovesFromPosition(req.UserID, position)
}

func newGameRequest(req *ws.ClientMessage[websocket.NewGameRequest]) (*Board, error) {
	var (
		whitePlayerID, blackPlayerID *uuid.UUID
	)

	if req.Data.Color == "white" {
		whitePlayerID = &req.UserID
	} else {
		blackPlayerID = &req.UserID
	}

	return newGame(req.Ctx, whitePlayerID, blackPlayerID)
}

func userBoardsRequset(req *ws.ClientMessage[any]) ([]*BoardOutputModel, error) {
	gameIds, err := gameService.GetGamesByUserID(req.Ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	boards := make([]*BoardOutputModel, len(gameIds))

	for i, game := range gameIds {
		board, err := getBoard(req.Ctx, game.ID.String())
		if err != nil {
			return nil, err
		}
		boards[i] = board.ToOutput()
	}

	return boards, nil
}

func placePieceRequest(req *ws.ClientMessage[websocket.GameMovePieceRequest]) (*Board, error) {
	board, err := getBoard(req.Ctx, req.Data.GameID)

	if err != nil {
		return nil, err
	}

	if board.Game.Status == models.GameStatusWaiting {
		return nil, ErrGameWaitingStatus
	}

	from, err := chessboard.GetPosition(req.Data.From)

	if err != nil {
		return nil, err
	}

	to, err := chessboard.GetPosition(req.Data.To)

	if err != nil {
		return nil, err
	}

	if err := board.PlacePieceFromPosition(req.Ctx, req.UserID, from, to); err != nil {
		return nil, err
	}

	return board, nil
}

func joinGameRequest(req *ws.ClientMessage[websocket.JoinGameRequest]) (*Board, error) {
	board, err := getBoard(req.Ctx, req.Data.GameID)

	if err != nil {
		return nil, err
	}

	if err := board.JoinPlayer(req.Ctx, req.UserID); err != nil {
		return nil, err
	}

	return board, nil
}
