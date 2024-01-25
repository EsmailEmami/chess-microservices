package websocket

import (
	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/websocket"
)

const (
	PrivateMsg   = "private-msg"
	Notification = "notification"
	ErrorMsg     = "error-msg"

	// game
	NewGame        = "new-game"
	GameValidMoves = "game-valid-moves"
	GameBoards     = "game-boards"
	GameMovePiece  = "game-move-piece"
	JoinGame       = "join-game"

	// send types
	GameBoardChanged = "game-board-changed"
)

var (
	NewGameCh        = make(chan *websocket.ClientMessage[NewGameRequest], 256)
	GameValidMovesCh = make(chan *websocket.ClientMessage[GameValidMovesRequest], 256)
	GameBoardsCh     = make(chan *websocket.ClientMessage[any], 256)
	GameMovePieceCh  = make(chan *websocket.ClientMessage[GameMovePieceRequest], 256)
	JoinGameCh       = make(chan *websocket.ClientMessage[JoinGameRequest], 256)
)

func ChessOnMessage(c *websocket.Client, msg *websocket.Message) {
	switch msg.Type {
	case NewGame:
		var req NewGameRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		NewGameCh <- websocket.NewClientMessage(c, req)

	case GameValidMoves:
		var req GameValidMovesRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		GameValidMovesCh <- websocket.NewClientMessage(c, req)

	case GameBoards:
		GameBoardsCh <- websocket.NewClientMessage[any](c, nil)
	case GameMovePiece:
		var req GameMovePieceRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		GameMovePieceCh <- websocket.NewClientMessage(c, req)
	case JoinGame:
		var req JoinGameRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}
		JoinGameCh <- websocket.NewClientMessage(c, req)
	default:
		logging.Warn("websocket invalid message type", msg.Type)
	}
}
