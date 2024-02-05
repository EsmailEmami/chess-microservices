package websocket

import (
	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/websocket"
)

const (
	// game
	ChessValidMoves = "chess-valid-moves"
	ChessMovePiece  = "chess-move-piece"

	// send types
	NewBoard          = "new-board"
	ChessInCheck      = "chess-in-check"
	ChessCheckmate    = "chess-chackmate"
	ChessPlayerJoined = "chess-player-joined"
	ChessNewWatcher   = "chess-new-watcher"
)

var (
	ChessRegisterCh   = make(chan *websocket.Client, 256)
	ChessUnregisterCh = make(chan *websocket.Client, 256)
	ChessValidMovesCh = make(chan *websocket.ClientMessage[ChessValidMovesRequest], 256)
	ChessMovePieceCh  = make(chan *websocket.ClientMessage[ChessMovePieceRequest], 256)
)

func ChessOnMessage(c *websocket.Client, msg *websocket.Message) {
	switch msg.Type {
	case ChessValidMoves:
		var req ChessValidMovesRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		ChessValidMovesCh <- websocket.NewClientMessage(c, req)
	case ChessMovePiece:
		var req ChessMovePieceRequest
		if !c.Unmarshal(msg.Content, &req) {
			return
		}

		ChessMovePieceCh <- websocket.NewClientMessage(c, req)
	default:
		logging.Warn("websocket invalid message type", "type", msg.Type)
	}
}

func ChessOnRegister(c *websocket.Client) {
	ChessRegisterCh <- c
}

func ChessOnUnregister(c *websocket.Client) {
	ChessUnregisterCh <- c
}
