package routes

import (
	apiHandler "github.com/esmailemami/chess/shared/handler"
	"github.com/esmailemami/chess/user/api/handler"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/gin-gonic/gin"
)

func friendRoutes(r *gin.RouterGroup, friendService *service.FriendService) {
	api := r.Group("/friend")

	handler := handler.NewFriendHandler(friendService)

	api.POST("/make/:friendId", apiHandler.HandleAPI(handler.MakeFriend))
	api.POST("/remove/:friendId", apiHandler.HandleAPI(handler.RemoveFriend))
	api.GET("/", apiHandler.HandleAPI(handler.GetFriends))
}
