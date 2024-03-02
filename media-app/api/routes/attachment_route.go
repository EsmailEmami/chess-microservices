package routes

import (
	"github.com/esmailemami/chess/media/api/handler"
	"github.com/esmailemami/chess/media/internal/app/service"

	apiHandler "github.com/esmailemami/chess/shared/handler"
	"github.com/gin-gonic/gin"
)

func attachmentRoutes(r *gin.RouterGroup, attachmentService *service.AttachmentService) {
	api := r.Group("/attachment")

	attachmentHandler := handler.NewAttachmentHandler(attachmentService)

	api.POST("/upload/profile/:id", apiHandler.HandleAPI(attachmentHandler.UploadProile))
	api.POST("/upload/room/avatar/:id", apiHandler.HandleAPI(attachmentHandler.UploadRoomAvatar))
	api.POST("/upload/room/file-message/:id", apiHandler.HandleAPI(attachmentHandler.UploadRoomFileMessage))

	api.GET("/download/stream/:id", apiHandler.HandleAPI(attachmentHandler.Stream))
}
