package routes

import (
	"github.com/esmailemami/chess/game/internal/app/service"
	"github.com/esmailemami/chess/shared/database/redis"
	sharedService "github.com/esmailemami/chess/shared/service"
	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")
	chessService := service.NewChessService(redis.GetConnection(), sharedService.NewUserService())

	chessRoutes(route, chessService)
}
