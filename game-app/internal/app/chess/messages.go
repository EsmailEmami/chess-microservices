package chess

import (
	"github.com/esmailemami/chess/game/internal/app/models"
	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/shared/websocket"
	"github.com/google/uuid"
)

type ChessMessage struct {
	ChessID uuid.UUID `json:"chessId"`
	Data    any       `json:"data"`
}

type MovePieceResponse struct {
	From *chessboard.Position `json:"from"`
	To   *chessboard.Position `json:"to"`
}

type ChessOutPutResponse struct {
	models.ChessOutputModel

	IsCheckmate bool `json:"isCheckmate"`
	IsInCheck   bool `json:"isInCheck"`

	Connections []*ChessConnectionOutPutModel `json:"connections"`
}

type ChessConnectionOutPutModel struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `gorm:"first_name" json:"firstName"`
	LastName  *string   `gorm:"last_name" json:"lastName"`
	Username  string    `gorm:"username" json:"username"`
}

func NewChessConnection(client *websocket.Client) *ChessConnectionOutPutModel {
	return &ChessConnectionOutPutModel{
		ID:        client.UserID,
		FirstName: client.User.FirstName,
		LastName:  client.User.LastName,
		Username:  client.User.Username,
	}
}
