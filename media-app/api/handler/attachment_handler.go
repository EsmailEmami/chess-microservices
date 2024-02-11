package handler

import (
	"context"

	"github.com/esmailemami/chess/media/internal/app/service"
	dbModels "github.com/esmailemami/chess/media/internal/models"
	"github.com/esmailemami/chess/media/internal/rabbitmq"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentHandler struct {
	handler.Handler

	attachmentService *service.AttachmentService
}

func NewAttachmentHandler(attachmentService *service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{
		attachmentService: attachmentService,
	}
}

// UploadProile godoc
// @Tags attachment
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Param file formData file true "Image file to be uploaded"
// @Success 200 {object} handler.Response[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/upload/profile/{id} [post]
func (a *AttachmentHandler) UploadProile(ctx *gin.Context, id uuid.UUID) (*handler.Response[uuid.UUID], error) {
	files, err := a.GetFiles(ctx, handler.TenMB)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errs.BadRequestErr().Msg("No file received!")
	}

	attachment, err := a.attachmentService.UploadFile(ctx, files[0], id, dbModels.ATTACHMENT_USER_PROFILE, dbModels.ATTACHMENT_USER_PROFILE)
	if err != nil {
		return nil, err
	}

	// send the data to rabbitmq

	rabbitmq.PublishUserProfile(context.Background(), id, attachment.UploadPath)

	return handler.OK(&attachment.ID), nil
}
