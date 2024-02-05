package chess

import (
	"context"

	"github.com/esmailemami/chess/game/internal/app/service"
	"github.com/esmailemami/chess/game/pkg/websocket"
	"github.com/esmailemami/chess/shared/database/redis"
	sharedService "github.com/esmailemami/chess/shared/service"
	sharedWebsocket "github.com/esmailemami/chess/shared/websocket"
	"github.com/google/uuid"
)

var chessService *service.ChessService

func Run() {
	chessService = service.NewChessService(redis.GetConnection(), sharedService.NewUserService())

	for {
		select {
		case req := <-websocket.ChessValidMovesCh:
			chessValidMovesRequest(req)

		case req := <-websocket.ChessMovePieceCh:
			chessMovePieceRequest(req)

		case client := <-websocket.ChessRegisterCh:
			clientOnRegister(client)

		case client := <-websocket.ChessUnregisterCh:
			clientOnUnregister(client)
		}
	}
}

func chessValidMovesRequest(req *sharedWebsocket.ClientMessage[websocket.ChessValidMovesRequest]) {
	board, err := getBoard(req.Ctx, req.Data.GameID)

	if err != nil {
		websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	data, err := board.GetValidMoves(req)

	if err != nil {
		websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	websocket.ChessWss.SendMessageToClient(req.ClientID, websocket.ChessValidMoves, &ChessMessage{
		ChessID: board.ChessID,
		Data:    data,
	})
}

func chessMovePieceRequest(req *sharedWebsocket.ClientMessage[websocket.ChessMovePieceRequest]) {
	board, err := getBoard(req.Ctx, req.Data.GameID)

	if err != nil {
		websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	from, to, err := board.PlacePiece(req)

	if err != nil {
		websocket.ChessWss.SendErrorMessageToClient(req.ClientID, err.Error())
		return
	}

	for _, client := range board.connections {
		websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.ChessMovePiece, &ChessMessage{
			ChessID: board.ChessID,
			Data: &MovePieceResponse{
				From: from,
				To:   to,
			},
		})
	}

	// check for check or checkmate
	if board.IsCheckmate() {
		output, err := board.OutPut()

		if err != nil {
			for _, client := range board.connections {
				websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.ChessCheckmate, &ChessMessage{
					ChessID: board.ChessID,
					Data:    output,
				})
			}
		}

		// its a checkmate, we have to delete the chess game from the map games
		deleteChess(board.ChessID)
	} else if board.IsInCheck() {
		output, err := board.OutPut()

		if err != nil {
			for _, client := range board.connections {
				websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.ChessInCheck, &ChessMessage{
					ChessID: board.ChessID,
					Data:    output,
				})
			}
		}
	}
}

func clientOnRegister(client *sharedWebsocket.Client) {
	chessIDs, err := chessService.GetChessIDsByUser(client.Context, client.UserID)

	if err != nil {
		websocket.ChessWss.SendErrorMessageToClient(client.SessionID, err.Error())
		return
	}

	for _, chessID := range chessIDs {
		board, err := getBoard(client.Context, chessID)
		if err != nil {
			websocket.ChessWss.SendErrorMessageToClient(client.SessionID, err.Error())
			continue
		}

		output, err := board.OutPut()
		if err != nil {
			websocket.ChessWss.SendErrorMessageToClient(client.SessionID, err.Error())
			continue
		}

		websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.NewBoard, &ChessMessage{
			ChessID: board.ChessID,
			Data:    output,
		})
	}
}

func clientOnUnregister(client *sharedWebsocket.Client) {
	chessIDs, err := chessService.GetChessIDsByUser(client.Context, client.UserID)

	if err != nil {
		websocket.ChessWss.SendErrorMessageToClient(client.SessionID, err.Error())
		return
	}

	chessIDs = append(chessIDs, chessService.GetUserWatchingGames(client.UserID)...)

	for _, chessID := range chessIDs {
		board, err := getBoard(client.Context, chessID)
		if err != nil {
			board.Disconnect(client)
		}

	}

	// delete watching cache
	chessService.DeleteWatcherCache(client.UserID)
}

func Join(ctx context.Context, userID, chessID uuid.UUID) error {
	board, err := getBoard(ctx, chessID)

	if err != nil {
		return err
	}

	for _, client := range websocket.ChessWss.GetUserConnections(userID) {
		board.Connect(client)
	}

	board.JoinPlayer(ctx, userID)

	output, err := board.OutPut()
	if err != nil {
		for _, client := range board.connections {
			websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.ChessPlayerJoined, &ChessMessage{
				ChessID: board.ChessID,
				Data:    output,
			})

		}
	}

	return nil
}

func Watch(ctx context.Context, userID, chessID uuid.UUID) error {
	board, err := getBoard(ctx, chessID)

	if err != nil {
		return err
	}

	newClients := websocket.ChessWss.GetUserConnections(userID)

	for _, client := range newClients {
		board.Connect(client)
	}

	output, err := board.OutPut()
	if err != nil {
		return err
	}

	for _, client := range board.connections {
		// if new client, send the whole board
		isNew := false

		for _, newClient := range newClients {
			if newClient.SessionID == client.SessionID {
				isNew = true
				break
			}
		}

		if isNew {
			websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.NewBoard, &ChessMessage{
				ChessID: board.ChessID,
				Data:    output,
			})
		} else {
			websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.ChessNewWatcher, &ChessMessage{
				ChessID: board.ChessID,
				Data:    NewChessConnection(client),
			})
		}

	}

	return nil
}

func New(ctx context.Context, userID, chessID uuid.UUID) error {
	board, err := getBoard(ctx, chessID)
	if err != nil {
		return err
	}

	output, err := board.OutPut()
	if err != nil {
		return err
	}

	// for two players

	if board.WhitePlayerUserID != nil {
		for _, client := range websocket.ChessWss.GetUserConnections(*board.WhitePlayerUserID) {
			board.Connect(client)

			websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.NewBoard, &ChessMessage{
				ChessID: board.ChessID,
				Data:    output,
			})
		}
	}

	if board.BlackPlayerUserID != nil {
		for _, client := range websocket.ChessWss.GetUserConnections(*board.BlackPlayerUserID) {
			board.Connect(client)

			websocket.ChessWss.SendMessageToClient(client.SessionID, websocket.NewBoard, &ChessMessage{
				ChessID: board.ChessID,
				Data:    output,
			})
		}
	}

	return nil
}
