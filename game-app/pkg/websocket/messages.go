package websocket

type NewGameRequest struct {
	Color string `json:"color"`
}

type GameValidMovesRequest struct {
	GameID   string `json:"gameId"`
	Position string `json:"position"`
}

type GameMovePieceRequest struct {
	GameID string `json:"gameId"`
	From   string `json:"position"`
	To     string `json:"to"`
}

type JoinGameRequest struct {
	GameID string `json:"gameId"`
}
