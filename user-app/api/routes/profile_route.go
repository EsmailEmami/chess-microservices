package routes

import (
	apiHandler "github.com/esmailemami/chess/shared/handler"
	"github.com/esmailemami/chess/user/api/handler"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/gin-gonic/gin"
)

func profileRoutes(r *gin.RouterGroup, userService *service.UserService) {
	api := r.Group("/profile")

	profileHandler := handler.NewProfileHandler(userService)

	api.GET("/", apiHandler.HandleAPI(profileHandler.Profile))
	api.POST("", apiHandler.HandleAPI(profileHandler.ChangeProfile))
	api.POST("/change-password", apiHandler.HandleAPI(profileHandler.ChangePassword))
}
