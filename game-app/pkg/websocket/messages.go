package websocket

import "github.com/google/uuid"

type ChessValidMovesRequest struct {
	GameID   uuid.UUID `json:"gameId"`
	Position string    `json:"position"`
}

type ChessMovePieceRequest struct {
	GameID uuid.UUID `json:"gameId"`
	From   string    `json:"position"`
	To     string    `json:"to"`
}
