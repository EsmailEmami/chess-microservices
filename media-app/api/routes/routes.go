package routes

import (
	"github.com/esmailemami/chess/media/internal/app/service"
	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {
	route := r.Group("api/v1")

	attachmentService := service.NewAttachmentService()
	attachmentRoutes(route, attachmentService)
}
