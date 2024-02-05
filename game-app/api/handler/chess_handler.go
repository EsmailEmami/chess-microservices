package handler

import (
	"github.com/esmailemami/chess/game/internal/app/chess"
	"github.com/esmailemami/chess/game/internal/app/models"
	"github.com/esmailemami/chess/game/internal/app/service"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/handler"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChessHandler struct {
	handler.Handler

	chessService *service.ChessService
}

func NewChessHandler(chessService *service.ChessService) *ChessHandler {
	return &ChessHandler{
		chessService: chessService,
	}
}

func (g *ChessHandler) JoinGame(ctx *gin.Context, id uuid.UUID) (*handler.Response[bool], error) {
	currentUser, err := g.GetUser(ctx)

	if err != nil {
		return nil, errs.UnAuthorizedErr()
	}

	err = g.chessService.JoinGame(ctx, currentUser, id)

	if err != nil {
		return nil, err
	}

	if err := chess.Join(ctx, currentUser.ID, id); err != nil {
		logging.WarnE("failed to join chess in websocket", err)
	}

	return handler.OKBool(), nil
}

func (g *ChessHandler) WatchGame(ctx *gin.Context, id uuid.UUID) (*handler.Response[bool], error) {
	currentUser, err := g.GetUser(ctx)

	if err != nil {
		return nil, errs.UnAuthorizedErr()
	}

	if err := chess.Watch(ctx, currentUser.ID, id); err != nil {
		return nil, err
	}

	return handler.OKBool(), nil
}

func (g *ChessHandler) NewChess(ctx *gin.Context, req models.CreateChessInputModel) (*handler.Response[bool], error) {
	currentUser, err := g.GetUser(ctx)

	if err != nil {
		return nil, errs.UnAuthorizedErr()
	}

	dbChess, err := g.chessService.NewChess(ctx, currentUser, &req)
	if err != nil {
		return nil, err
	}

	// send the new chess to the websocket to show users games
	if err := chess.New(ctx, currentUser.ID, dbChess.ID); err != nil {
		logging.WarnE("failed to create chess in websocket", err)
	}

	return handler.OKBool(), nil
}
