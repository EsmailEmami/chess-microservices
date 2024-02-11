package routes

import (
	"github.com/esmailemami/chess/media/internal/app/service"
	"github.com/esmailemami/chess/shared/middleware"
	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")
	route.Use(middleware.Authorization())

	attachmentService := service.NewAttachmentService()
	attachmentRoutes(route, attachmentService)
}
