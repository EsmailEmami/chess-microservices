package routes

import (
	"github.com/esmailemami/chess/shared/middleware"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")
	route.Use(middleware.Authorization())

	userService := service.NewUserService()

	profileRoutes(route, userService)
}
