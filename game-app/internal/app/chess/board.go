package chess

import (
	"context"
	"sync"

	"github.com/esmailemami/chess/game/internal/app/service"
	"github.com/esmailemami/chess/game/internal/models"
	"github.com/esmailemami/chess/game/pkg/chessboard"
	"github.com/esmailemami/chess/game/pkg/websocket"
	sharedWebsocket "github.com/esmailemami/chess/shared/websocket"
	"github.com/google/uuid"
)

type Board struct {
	chess             *chessboard.Chessboard
	ChessID           uuid.UUID
	WhitePlayerUserID *uuid.UUID
	BlackPlayerUserID *uuid.UUID
	Turn              uuid.UUID
	Status            models.ChessStatus
	chessService      *service.ChessService

	mutex sync.Mutex

	connections map[uuid.UUID]*sharedWebsocket.Client
}

func newBoard(chessID uuid.UUID, whitePlayerID, blackPlayerID *uuid.UUID, status models.ChessStatus, chessboard *chessboard.Chessboard, chessService *service.ChessService) (*Board, error) {
	board := &Board{
		mutex:             sync.Mutex{},
		chess:             chessboard,
		ChessID:           chessID,
		WhitePlayerUserID: whitePlayerID,
		BlackPlayerUserID: blackPlayerID,
		chessService:      chessService,
		Status:            status,
		connections:       make(map[uuid.UUID]*sharedWebsocket.Client),
	}

	board.setTurn()

	return board, nil
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

func (c *Board) isValidUser(userID uuid.UUID) bool {
	return userID == *c.WhitePlayerUserID || userID == *c.BlackPlayerUserID
}

func (c *Board) isValidTurn(userID uuid.UUID) bool {
	return c.Turn == userID
}

func (c *Board) getOppponentID() uuid.UUID {
	if c.Turn == *c.WhitePlayerUserID {
		return *c.BlackPlayerUserID
	} else {
		return *c.WhitePlayerUserID
	}
}

func (c *Board) swichTurn() {
	c.Turn = c.getOppponentID()
}

func (b *Board) setTurn() {
	if b.WhitePlayerUserID != nil {
		b.Turn = *b.WhitePlayerUserID
	} else {
		b.Turn = *b.BlackPlayerUserID
	}
}

func (b *Board) IsInCheck() bool {
	return b.chess.IsInCheck(b.getTurnColor())
}

func (b *Board) IsCheckmate() bool {
	return b.chess.IsCheckmate(b.getTurnColor())
}

func (b *Board) getTurnColor() chessboard.Color {
	if b.Turn == *b.WhitePlayerUserID {
		return chessboard.White
	}
	return chessboard.Black
}

func (b *Board) JoinPlayer(ctx context.Context, playerUserID uuid.UUID) error {
	if b.Status != models.ChessStatusWaiting {
		return ErrGameIsNotInWaitingStatus
	}

	if b.WhitePlayerUserID == nil {
		b.WhitePlayerUserID = &playerUserID

	} else {
		b.BlackPlayerUserID = &playerUserID
	}

	return nil
}

func (b *Board) OutPut() (*ChessOutPutResponse, error) {
	chess, err := b.chessService.Get(context.Background(), b.ChessID)
	if err != nil {
		return nil, err
	}

	//TODO: needs to fix this per user, not connections, may be one user with two connections
	connections := make([]*ChessConnectionOutPutModel, len(b.connections))

	it := 0
	for _, client := range b.connections {
		connections[it] = NewChessConnection(client)
		it++
	}

	output := &ChessOutPutResponse{
		IsInCheck:        b.IsInCheck(),
		IsCheckmate:      b.IsCheckmate(),
		ChessOutputModel: *chess,
	}

	return output, nil
}

// websocket calls

func (b *Board) GetValidMoves(req *sharedWebsocket.ClientMessage[websocket.ChessValidMovesRequest]) ([]chessboard.Position, error) {
	if err := b.userchecks(req.UserID); err != nil {
		return nil, err
	}

	pos, err := chessboard.GetPosition(req.Data.Position)

	if err != nil {
		return nil, err
	}

	validMoves := b.chess.GetValidMovesFromPosition(*pos)

	return validMoves, nil
}

func (b *Board) PlacePiece(req *sharedWebsocket.ClientMessage[websocket.ChessMovePieceRequest]) (from, to *chessboard.Position, err error) {
	if err = b.userchecks(req.UserID); err != nil {
		return nil, nil, err
	}

	if b.Status == models.ChessStatusWaiting {
		return nil, nil, ErrGameWaitingStatus
	}

	if b.Status != models.ChessStatusOpen {
		return nil, nil, ErrGameIsOver
	}

	from, err = chessboard.GetPosition(req.Data.From)

	if err != nil {
		return nil, nil, err
	}

	to, err = chessboard.GetPosition(req.Data.To)

	if err != nil {
		return nil, nil, err
	}

	piece := b.chess.GetPiece(from.Row, from.Col)
	if err := b.chess.PlacePiece(piece, *to); err != nil {
		return nil, nil, err
	}

	b.swichTurn()

	if err := b.chessService.MoveChessPiece(req.Ctx, b.ChessID, piece, *from, *to); err != nil {
		return nil, nil, err
	}

	if b.IsCheckmate() {
		if err := b.chessService.Chectmate(req.Ctx, b.ChessID, b.Turn); err != nil {
			return nil, nil, err
		}

		// the game is end, so it is close
		b.Status = models.ChessStatusClose
	}

	return from, to, nil
}

func (b *Board) Connect(client *sharedWebsocket.Client) {
	b.mutex.Lock()
	b.connections[client.SessionID] = client
	b.mutex.Unlock()
}

func (b *Board) Disconnect(client *sharedWebsocket.Client) {
	b.mutex.Lock()
	delete(b.connections, client.SessionID)
	b.mutex.Unlock()
}
