package routes

import (
	"github.com/esmailemami/chess/game/internal/app/service"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/middleware"
	sharedService "github.com/esmailemami/chess/shared/service"
	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")
	route.Use(middleware.Authorization())

	chessService := service.NewChessService(redis.GetConnection(), sharedService.NewUserService())

	chessRoutes(route, chessService)
}
