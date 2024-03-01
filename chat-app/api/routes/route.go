package routes

import (
	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/middleware"
	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")
	route.Use(middleware.Authorization())

	var (
		cache          = redis.GetConnection()
		messageService = service.NewMessageService(cache)
		roomService    = service.NewRoomService(cache, messageService)
	)

	roomRoutes(route, roomService)
}
