package models

import (
	"github.com/esmailemami/chess/game/internal/models"
	"github.com/google/uuid"
)

type ChessOutputModel struct {
	ID            uuid.UUID               `json:"id"`
	WhitePlayerID *uuid.UUID              `json:"whitePlayerId"`
	WhitePlayer   *ChessPlayerOutputModel `json:"whitePlayer"`
	BlackPlayerID *uuid.UUID              `json:"blackPlayerId"`
	BlackPlayer   *ChessPlayerOutputModel `json:"blackPlayer"`
	Turn          models.ChessPlayer      `json:"turn"`
	Moves         models.ChessMoves       `json:"moves"`
	Pieces        models.ChessPieces      `json:"pieces"`
	Status        models.ChessStatus      `json:"status"`
	IsInCheck     bool                    `json:"isCheck"`
	IsCheckmate   bool                    `json:"isCheckmate"`
	Winner        *uuid.UUID              `json:"winner"`
}

type ChessPlayerOutputModel struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `gorm:"first_name" json:"firstName"`
	LastName  *string   `gorm:"last_name" json:"lastName"`
	Username  string    `gorm:"username" json:"username"`
}

type CreateChessInputModel struct {
	Color       string     `json:"color"`
	PlayingWith *uuid.UUID `json:"playingWith,omitempty"`
}
