package game

import (
	"context"

	"github.com/esmailemami/chess/game/internal/models"
	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/shared/util"
	"github.com/google/uuid"
)

type Board struct {
	chess             *chessboard.Board
	Game              *models.Game
	WhitePlayerUserID *uuid.UUID
	BlackPlayerUserID *uuid.UUID
	Turn              uuid.UUID
}

func newBoard(g *models.Game, chessboard *chessboard.Board) (*Board, error) {
	board := &Board{
		chess:             chessboard,
		Game:              g,
		WhitePlayerUserID: g.WhitePlayerID,
		BlackPlayerUserID: g.BlackPlayerID,
	}

	board.setTurn()

	return board, nil
}

func (c *Board) PlacePieceFromPosition(ctx context.Context, userID uuid.UUID, from, to chessboard.Position) error {
	if err := c.userchecks(userID); err != nil {
		return err
	}

	piece := c.chess.GetPiece(from.Row, from.Col)

	// place the piece in the chessboard
	if err := c.chess.PlacePiece(piece, to); err != nil {
		return err
	}

	// remove from and to positions from the game
	for i, piece := range c.Game.Pieces {
		if !(piece.Row == from.Row && piece.Col == from.Col) && !(piece.Row == to.Row && piece.Col == to.Col) {
			continue
		}

		c.Game.Pieces = util.ArrayRemoveIndex[models.GamePiece](c.Game.Pieces, i)
	}

	// save the new piece position and move
	player := models.GetGamePlayerFromChessColor(piece.Color)

	c.Game.Pieces = append(c.Game.Pieces, models.GamePiece{
		Piece:  string(piece.Type),
		Row:    to.Row,
		Col:    to.Col,
		Player: player,
	})

	c.Game.Moves = append(c.Game.Moves, models.GameMove{
		Player: player,
		From:   from.String(),
		To:     to.String(),
	})

	if err := gameService.Update(ctx, c.Game); err != nil {
		return err
	}

	c.swichTurn()
	return nil
}

func (c *Board) GetValidMovesFromPosition(userID uuid.UUID, position chessboard.Position) ([]chessboard.Position, error) {
	if err := c.userchecks(userID); err != nil {
		return nil, err
	}

	return c.chess.GetValidMovesFromPosition(position), nil
}

func (c *Board) isValidUser(userID uuid.UUID) bool {
	return userID == *c.WhitePlayerUserID || userID == *c.BlackPlayerUserID
}

func (c *Board) isValidTurn(userID uuid.UUID) bool {
	return c.Turn == userID
}

func (c *Board) userchecks(userID uuid.UUID) error {
	if c.BlackPlayerUserID == nil || c.WhitePlayerUserID == nil {
		return ErrGameWaitingStatus
	}

	if !c.isValidUser(userID) {
		return ErrInvalidGame
	}

	if !c.isValidTurn(userID) {
		return ErrInvalidTurn
	}
	return nil
}

func (c *Board) swichTurn() {
	if c.Turn == *c.WhitePlayerUserID {
		c.Turn = *c.BlackPlayerUserID
	} else {
		c.Turn = *c.WhitePlayerUserID
	}
}

func (b *Board) setTurn() {
	if b.Game.Turn == models.GamePlayerWhite {
		b.Turn = *b.WhitePlayerUserID
	} else {
		b.Turn = *b.BlackPlayerUserID
	}
}

func (b *Board) ToOutput() *BoardOutputModel {
	o := &BoardOutputModel{
		GameID: b.Game.ID,

		Turn:   b.Turn,
		Moves:  b.Game.Moves,
		Pieces: b.Game.Pieces,
		Status: b.Game.Status,
	}

	if b.Game.WhitePlayer != nil {
		o.WhitePlayer = &BoardPlayerOutputModel{
			ID:        *b.WhitePlayerUserID,
			FirstName: b.Game.WhitePlayer.FirstName,
			LastName:  b.Game.WhitePlayer.LastName,
			Username:  b.Game.WhitePlayer.Username,
		}
	}

	if b.Game.BlackPlayer != nil {
		o.BlackPlayer = &BoardPlayerOutputModel{
			ID:        *b.BlackPlayerUserID,
			FirstName: b.Game.BlackPlayer.FirstName,
			LastName:  b.Game.BlackPlayer.LastName,
			Username:  b.Game.BlackPlayer.Username,
		}
	}

	return o
}

func (b *Board) JoinPlayer(ctx context.Context, playerUserID uuid.UUID) error {
	if b.Game.Status != models.GameStatusWaiting {
		return ErrGameIsNotInWaitingStatus
	}

	if b.WhitePlayerUserID == nil {
		b.WhitePlayerUserID = &playerUserID
		b.Game.WhitePlayerID = &playerUserID

		user, err := userService.Get(ctx, playerUserID)

		if err != nil {
			return err
		}

		b.Game.WhitePlayer = user

	} else {
		b.BlackPlayerUserID = &playerUserID
		b.Game.BlackPlayerID = &playerUserID

		user, err := userService.Get(ctx, playerUserID)

		if err != nil {
			return err
		}

		b.Game.BlackPlayer = user
	}

	b.Game.Status = models.GameStatusOpen

	if err := gameService.Update(ctx, b.Game); err != nil {
		return err
	}

	return nil
}

type BoardOutputModel struct {
	GameID      uuid.UUID               `json:"gameId"`
	WhitePlayer *BoardPlayerOutputModel `json:"whitePlayer"`
	BlackPlayer *BoardPlayerOutputModel `json:"blackPlayer"`
	Turn        uuid.UUID               `json:"turn"`
	Moves       models.GameMoves        `json:"moves"`
	Pieces      models.GamePieces       `json:"pieces"`
	Status      models.GameStatus       `json:"status"`
}

type BoardPlayerOutputModel struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `gorm:"first_name" json:"firstName"`
	LastName  *string   `gorm:"last_name" json:"lastName"`
	Username  string    `gorm:"username" json:"username"`
}
