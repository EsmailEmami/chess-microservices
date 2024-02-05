package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

type ChessStatus int

const (
	ChessStatusWaiting ChessStatus = iota
	ChessStatusRejected
	ChessStatusOpen
	ChessStatusClose
)

type ChessPlayer string

const (
	ChessPlayerWhite = "white"
	ChessPlayerBlack = "black"
)

func GetChessPlayerFromColor(color chessboard.Color) ChessPlayer {
	if color == chessboard.Black {
		return ChessPlayerBlack
	} else {
		return ChessPlayerWhite
	}
}

func (g ChessPlayer) ChessColor() chessboard.Color {
	if g == ChessPlayerBlack {
		return chessboard.Black
	} else {
		return chessboard.White
	}
}

type Chess struct {
	models.Model

	WhitePlayerID *uuid.UUID   `gorm:"white_player_id" json:"whitePlayerId"`
	WhitePlayer   *models.User `gorm:"foreignKey:white_player_id;references:id" json:"whitePlayer"`
	BlackPlayerID *uuid.UUID   `gorm:"black_player_id" json:"blackPlayerId"`
	BlackPlayer   *models.User `gorm:"foreignKey:black_player_id;references:id" json:"blackPlayer"`
	Turn          ChessPlayer  `gorm:"turn" json:"turn"`
	Moves         ChessMoves   `gorm:"moves" json:"moves"`
	Pieces        ChessPieces  `gorm:"pieces" json:"pieces"`
	Status        ChessStatus  `gorm:"status" json:"status"`
	WinnerID      *uuid.UUID   `gorm:"winner_id" json:"winnerId"`
	Winner        *models.User `gorm:"foreignKey:winner_id;references:id" json:"winner"`
}

func (Chess) TableName() string {
	return "game.chess"
}

func NewChess(whitePlayer, blackPlayer *models.User, pieces []*chessboard.ChessboardPiece) *Chess {
	chess := &Chess{
		WhitePlayer: whitePlayer,
		BlackPlayer: blackPlayer,
		Status:      ChessStatusWaiting,
		Turn:        ChessPlayerWhite,
		Pieces:      make(ChessPieces, len(pieces)),
		Moves:       make(ChessMoves, 0),
	}
	chess.ID = uuid.New()

	if whitePlayer != nil {
		chess.WhitePlayerID = &whitePlayer.ID
		chess.Turn = ChessPlayerWhite
	}

	if blackPlayer != nil {
		chess.BlackPlayerID = &blackPlayer.ID
		chess.Turn = ChessPlayerBlack
	}

	// the Chess is accepted and open
	if whitePlayer != nil && blackPlayer != nil {
		chess.Status = ChessStatusOpen
	}

	// set default pieces
	for i, piece := range pieces {
		chess.Pieces[i] = ChessPiece{
			Row:    piece.Row,
			Col:    piece.Col,
			Piece:  string(piece.PieceType),
			Player: GetChessPlayerFromColor(piece.Color),
		}
	}

	return chess
}

func (g *Chess) SwitchTurn() {
	if g.Turn == ChessPlayerWhite {
		g.Turn = ChessPlayerBlack
	} else {
		g.Turn = ChessPlayerWhite
	}
}

// ChessMoves
type ChessMove struct {
	Player ChessPlayer
	From   string // Example: a1
	To     string // Example: a2
}

type ChessMoves []ChessMove

func (p ChessMoves) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *ChessMoves) Scan(value interface{}) error {
	if value == nil {
		j = nil
		return nil
	}
	var bts []byte
	switch v := value.(type) {
	case []byte:
		bts = v
	case string:
		bts = []byte(v)
	case nil:
		*j = nil
		return nil
	}
	return json.Unmarshal(bts, &j)
}

// ChessPiece

type ChessPiece struct {
	Player   ChessPlayer
	Piece    string
	Row, Col int
}

func (p *ChessPiece) ToChessPiece() *chessboard.ChessboardPiece {
	return &chessboard.ChessboardPiece{
		Row:       p.Row,
		Col:       p.Col,
		PieceType: chessboard.PieceType(p.Piece),
		Color:     p.Player.ChessColor(),
	}
}

type ChessPieces []ChessPiece

func (p ChessPieces) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *ChessPieces) Scan(value interface{}) error {
	if value == nil {
		j = nil
		return nil
	}
	var bts []byte
	switch v := value.(type) {
	case []byte:
		bts = v
	case string:
		bts = []byte(v)
	case nil:
		*j = nil
		return nil
	}
	return json.Unmarshal(bts, &j)
}
