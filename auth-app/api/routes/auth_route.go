package routes

import (
	"github.com/esmailemami/chess/auth/api/handler"
	"github.com/esmailemami/chess/auth/internal/app/service"
	apihandler "github.com/esmailemami/chess/shared/handler"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")

	var (
		userService = service.NewUserService()
		authService = service.NewAuthService(userService)
	)

	authHandler := handler.NewAuthHandler(authService)

	anonapi := route.Group("")
	anonapi.POST("/login", apihandler.HandleAPI(authHandler.Login))
	anonapi.POST("/register", apihandler.HandleAPI(authHandler.Register))
}
