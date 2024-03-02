package handler

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

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

// Stream godoc
// @Tags attachment
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Success 200 {object} handler.JSONResponse[bool]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/download/stream/{id} [GET]
func (a *AttachmentHandler) Stream(ctx *gin.Context, id uuid.UUID) (handler.Response, error) {
	filePath, mimeType, err := a.attachmentService.GetFileInfo(ctx, id)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, errs.NotFoundErr().Msg("file not found").WithError(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errs.NotFoundErr().Msg("file not found").WithError(err)
	}
	defer file.Close()

	ctx.Writer.Header().Set("Content-Type", mimeType)
	ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	bufferedWriter := bufio.NewWriter(ctx.Writer)
	defer bufferedWriter.Flush()

	const chunkSize = 10 * 1024
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		_, err = bufferedWriter.Write(buffer[:n])
		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		return nil, errs.InternalServerErr().Msg("Failed to stream file content to response").WithError(err)
	}

	return nil, nil
}

// UploadProile godoc
// @Tags attachment
// @Accept json
// @Produce json
// @Security Bearer
// @Param id   path  string  true  "id"
// @Param file formData file true "Image file to be uploaded"
// @Success 200 {object} handler.JSONResponse[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/upload/profile/{id} [post]
func (a *AttachmentHandler) UploadProile(ctx *gin.Context, id uuid.UUID) (handler.Response, error) {
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
// @Success 200 {object} handler.JSONResponse[uuid.UUID]
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/upload/room/avatar/{id} [post]
func (a *AttachmentHandler) UploadRoomAvatar(ctx *gin.Context, id uuid.UUID) (handler.Response, error) {
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
// @Failure 400 {object} errs.Error
// @Failure 422 {object} errs.ValidationError
// @Router /attachment/upload/room/file-message/{id} [post]
func (a *AttachmentHandler) UploadRoomFileMessage(ctx *gin.Context, id uuid.UUID) (handler.Response, error) {
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

	rabbitmq.PublishRoomFileMessage(context.Background(), id, messageID, user.ID, attachment.ID.String(), fileType)

	return nil, nil
}
