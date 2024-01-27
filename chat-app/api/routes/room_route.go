package routes

import (
	"github.com/esmailemami/chess/chat/api/handler"
	"github.com/esmailemami/chess/chat/internal/app/service"
	apiHandler "github.com/esmailemami/chess/shared/handler"

	"github.com/gin-gonic/gin"
)

func roomRoutes(r *gin.RouterGroup, roomService *service.RoomService) {
	api := r.Group("/room")

	roomHandler := handler.NewRoomHandler(roomService)

	api.POST("/private", apiHandler.HandleAPI(roomHandler.CreatePrivateRoom))
	api.POST("/global", apiHandler.HandleAPI(roomHandler.CreateGlobalRoom))
	api.GET("/:id", apiHandler.HandleAPI(roomHandler.GetRoom))
}
