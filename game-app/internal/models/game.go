package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

type GameStatus int

const (
	GameStatusWaiting GameStatus = iota
	GameStatusRejected
	GameStatusOpen
	GameStatusClose
)

type GamePlayer string

const (
	GamePlayerWhite = "white"
	GamePlayerBlack = "black"
)

func GetGamePlayerFromChessColor(color chessboard.Color) GamePlayer {
	if color == chessboard.Black {
		return GamePlayerBlack
	} else {
		return GamePlayerWhite
	}
}

func (g GamePlayer) ChessColor() chessboard.Color {
	if g == GamePlayerBlack {
		return chessboard.Black
	} else {
		return chessboard.White
	}
}

type Game struct {
	models.Model

	WhitePlayerID *uuid.UUID   `gorm:"white_player_id" json:"whitePlayerId"`
	WhitePlayer   *models.User `gorm:"foreignKey:white_player_id;references:id" json:"whitePlayer"`
	BlackPlayerID *uuid.UUID   `gorm:"black_player_id" json:"blackPlayerId"`
	BlackPlayer   *models.User `gorm:"foreignKey:black_player_id;references:id" json:"blackPlayer"`
	Turn          GamePlayer   `gorm:"turn" json:"turn"`
	Moves         GameMoves    `gorm:"moves" json:"moves"`
	Pieces        GamePieces   `gorm:"pieces" json:"pieces"`
	Status        GameStatus   `gorm:"status" json:"status"`
}

func (Game) TableName() string {
	return "game"
}

func NewGame(whitePlayer, blackPlayer *models.User, pieces []*chessboard.BoardPiece) *Game {
	game := &Game{
		WhitePlayer: whitePlayer,
		BlackPlayer: blackPlayer,
		Status:      GameStatusWaiting,
		Turn:        GamePlayerWhite,
		Pieces:      make(GamePieces, len(pieces)),
		Moves:       make(GameMoves, 0),
	}
	game.ID = uuid.New()

	if whitePlayer != nil {
		game.WhitePlayerID = &whitePlayer.ID
		game.Turn = GamePlayerWhite
	}

	if blackPlayer != nil {
		game.BlackPlayerID = &blackPlayer.ID
		game.Turn = GamePlayerBlack
	}

	// the game is accepted and open
	if whitePlayer != nil && blackPlayer != nil {
		game.Status = GameStatusOpen
	}

	// set default pieces
	for i, piece := range pieces {
		game.Pieces[i] = GamePiece{
			Row:    piece.Row,
			Col:    piece.Col,
			Piece:  string(piece.PieceType),
			Player: GetGamePlayerFromChessColor(piece.Color),
		}
	}

	return game
}

func (g *Game) GetChessPieces() []*chessboard.BoardPiece {
	pieces := make([]*chessboard.BoardPiece, len(g.Pieces))
	for i, piece := range g.Pieces {
		pieces[i] = piece.ToChessPiece()
	}

	return pieces
}

// GameMoves
type GameMove struct {
	Player GamePlayer
	From   string // Example: a1
	To     string // Example: a2
}

type GameMoves []GameMove

func (p GameMoves) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *GameMoves) Scan(value interface{}) error {
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

// GamePiece

type GamePiece struct {
	Player   GamePlayer
	Piece    string
	Row, Col int
}

func (p *GamePiece) ToChessPiece() *chessboard.BoardPiece {
	return &chessboard.BoardPiece{
		Row:       p.Row,
		Col:       p.Col,
		PieceType: chessboard.PieceType(p.Piece),
		Color:     p.Player.ChessColor(),
	}
}

type GamePieces []GamePiece

func (p GamePieces) Value() (driver.Value, error) {
	valueString, err := json.Marshal(p)
	return string(valueString), err
}

func (j *GamePieces) Scan(value interface{}) error {
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
