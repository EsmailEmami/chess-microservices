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

	api.GET("/", apiHandler.HandleAPI(roomHandler.GetRooms))
	api.POST("/private", apiHandler.HandleAPI(roomHandler.CreatePrivateRoom))
	api.POST("/public", apiHandler.HandleAPI(roomHandler.CreatePublicRoom))
	api.POST("/join/:id", apiHandler.HandleAPI(roomHandler.JoinRoom))
	api.POST("/left/:id", apiHandler.HandleAPI(roomHandler.LeftRoom))
}
