package routes

import (
	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/shared/middleware"
	sharedService "github.com/esmailemami/chess/shared/service"
	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")
	route.Use(middleware.Authorization())

	userService := sharedService.NewUserService()
	roomService := service.NewRoomService(userService)

	roomRoutes(route, roomService)
}
