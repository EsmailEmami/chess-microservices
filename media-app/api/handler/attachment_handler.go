package handler

import (
	"context"

	"github.com/esmailemami/chess/media/internal/app/service"
	dbModels "github.com/esmailemami/chess/media/internal/models"
	"github.com/esmailemami/chess/media/internal/util"
	"github.com/esmailemami/chess/media/pkg/rabbitmq"
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

	currentAttachmentID := a.attachmentService.GetCurrentAttachmentID(ctx, id)

	attachment, err := a.attachmentService.UploadFile(ctx, files[0], id, dbModels.ATTACHMENT_USER_PROFILE, dbModels.ATTACHMENT_USER_PROFILE)
	if err != nil {
		return nil, err
	}

	rabbitmq.PublishUserProfile(context.Background(), attachment.ID, id, attachment.UploadPath, currentAttachmentID)

	return handler.OK(&attachment.ID), nil
}

// UploadRoomAvatar godoc
// @Tags attachment
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Param file formData file true "Image file to be uploaded"
// @Success 200 {object} handler.Response[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/upload/room/avatar/{id} [post]
func (a *AttachmentHandler) UploadRoomAvatar(ctx *gin.Context, id uuid.UUID) (*handler.Response[uuid.UUID], error) {
	files, err := a.GetFiles(ctx, handler.TenMB)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errs.BadRequestErr().Msg("No file received!")
	}

	currentAttachmentID := a.attachmentService.GetCurrentAttachmentID(ctx, id)

	attachment, err := a.attachmentService.UploadFile(ctx, files[0], id, dbModels.ATTACHMENT_PUBLIC_ROOM_PROFILE, dbModels.ATTACHMENT_USER_PROFILE)
	if err != nil {
		return nil, err
	}

	rabbitmq.PublishRoomAvatar(context.Background(), attachment.ID, id, attachment.UploadPath, currentAttachmentID)

	return handler.OK(&attachment.ID), nil
}

// UploadRoomFileMessage godoc
// @Tags attachment
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Param file formData file true "Image file to be uploaded"
// @Success 200 {object} handler.Response[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/upload/room/file-message/{id} [post]
func (a *AttachmentHandler) UploadRoomFileMessage(ctx *gin.Context, id uuid.UUID) (*handler.Response[uuid.UUID], error) {
	user := a.GetUser(ctx)
	if user == nil {
		return nil, errs.UnAuthorizedErr()
	}

	files, err := a.GetFiles(ctx, handler.TenMB)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errs.BadRequestErr().Msg("No file received!")
	}

	file := files[0]
	mimeType := util.GetMimeType(file)

	if !util.IsImage(mimeType) && !util.IsVideo(mimeType) {
		return nil, errs.BadRequestErr().Msg("only image or video files are acceptable")
	}

	messageID := uuid.New()

	attachment, err := a.attachmentService.UploadFile(ctx, file, messageID, dbModels.ATTACHMENT_ROOM_FILE_MESSAGE, dbModels.ATTACHMENT_ROOM_FILE_MESSAGE)
	if err != nil {
		return nil, err
	}

	fileType := "image"
	if util.IsVideo(mimeType) {
		fileType = "video"
	}

	rabbitmq.PublishRoomFileMessage(context.Background(), id, messageID, user.ID, attachment.UploadPath, fileType)

	return handler.OK(&attachment.ID), nil
}
