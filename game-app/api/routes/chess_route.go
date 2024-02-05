package routes

import (
	"github.com/esmailemami/chess/game/api/handler"
	"github.com/esmailemami/chess/game/internal/app/service"
	apiHandler "github.com/esmailemami/chess/shared/handler"
	"github.com/gin-gonic/gin"
)

func chessRoutes(r *gin.RouterGroup, chessService *service.ChessService) {
	api := r.Group("/chess")

	roomHandler := handler.NewChessHandler(chessService)

	api.POST("/watch", apiHandler.HandleAPI(roomHandler.WatchGame))
	api.POST("/join", apiHandler.HandleAPI(roomHandler.JoinGame))
	api.POST("/", apiHandler.HandleAPI(roomHandler.JoinGame))
}
