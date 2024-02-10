package handler

import (
	"mime/multipart"

	"github.com/esmailemami/chess/media/internal/app/service"
	dbModels "github.com/esmailemami/chess/media/internal/models"
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
// @Success 200 {object} handler.Response[bool]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/upload/profile/{id} [post]
func (a *AttachmentHandler) UploadProile(ctx *gin.Context, id uuid.UUID) (*handler.Response[bool], error) {
	err := a.upload(ctx, id, dbModels.ATTACHMENT_USER_PROFILE)
	if err != nil {
		return nil, err
	}

	return handler.OKBool(), nil
}

func (a *AttachmentHandler) upload(c *gin.Context, itemID uuid.UUID, fileType string) error {
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB maximum file size
	if err != nil {
		return errs.BadRequestErr().WithError(err).Msg("File size is more than 10MB")
	}

	var files []*multipart.FileHeader

	for _, fileHeaders := range c.Request.MultipartForm.File {
		files = append(files, fileHeaders...)
	}

	if _, err := a.attachmentService.UploadFiles(c, files, itemID, fileType, fileType); err != nil {
		return err
	}

	return nil
}
