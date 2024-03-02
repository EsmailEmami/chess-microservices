package service

import (
	"context"
	"mime/multipart"

	"github.com/esmailemami/chess/media/internal/models"
	"github.com/esmailemami/chess/media/internal/util"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentService struct {
	fileService *FileService
}

func NewAttachmentService() *AttachmentService {
	return &AttachmentService{
		fileService: NewFileService(),
	}
}

func (a *AttachmentService) UploadFiles(ctx context.Context, files []*multipart.FileHeader, itemID uuid.UUID, directory string, fileType string) ([]*models.Attachment, error) {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	attachments := make([]*models.Attachment, 0, len(files))

	for _, file := range files {
		attachment, err := a.uploadFileTransaction(db, file, itemID, directory, fileType)
		if err != nil {
			tx.Rollback()

			for _, attachment := range attachments {
				_ = a.fileService.DeleteFile(attachment.UploadPath)
			}

			return nil, err
		}

		attachments = append(attachments, attachment)
	}

	tx.Commit()

	return attachments, nil
}

func (a *AttachmentService) UploadFile(ctx context.Context, f *multipart.FileHeader, itemID uuid.UUID, directory string, fileType string) (*models.Attachment, error) {
	db := psql.DBContext(ctx)
	return a.uploadFileTransaction(db, f, itemID, directory, fileType)
}

func (a *AttachmentService) uploadFileTransaction(db *gorm.DB, f *multipart.FileHeader, itemID uuid.UUID, directory string, fileType string) (*models.Attachment, error) {
	uploadPath, fileName, err := a.fileService.UploadFile(f, directory)

	if err != nil {
		return nil, err
	}

	attachment := &models.Attachment{
		MimeType:     util.GetMimeType(f),
		Extension:    util.GetFileExetension(fileName),
		UploadPath:   uploadPath,
		FileName:     fileName,
		OriginalName: f.Filename,
		FileType:     fileType,
		ItemID:       itemID,
	}

	if err := db.Create(attachment).Error; err != nil {
		// remove file
		_ = a.fileService.DeleteFile(uploadPath)

		return nil, errs.InternalServerErr().WithError(err)
	}

	return attachment, err
}

func (a *AttachmentService) Delete(ctx context.Context, attachmentID uuid.UUID) error {
	db := psql.DBContext(ctx)

	var attachment models.Attachment

	if err := db.Model(&models.Attachment{}).First(&attachment, "id = ?", attachmentID).Error; err != nil {
		return errs.NotFoundErr().WithError(err)
	}

	if err := a.fileService.DeleteFile(attachment.UploadPath); err != nil {
		return errs.InternalServerErr().Msg("failed to delete the file")
	}

	if err := db.Delete(&attachment).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (a *AttachmentService) GetCurrentAttachmentID(ctx context.Context, itemID uuid.UUID) *uuid.UUID {
	db := psql.DBContext(ctx)

	var attachmentID string

	if err := db.Model(&models.Attachment{}).Order("created_at DESC").Select("id").First(&attachmentID, "item_id = ?", itemID).Error; err != nil {
		logging.ErrorE("failed to get current attachment", err, "itemId", itemID)

		return nil
	}

	return func() *uuid.UUID {
		id := uuid.MustParse(attachmentID)
		return &id
	}()
}

func (a *AttachmentService) GetFileInfo(ctx context.Context, attachmentID uuid.UUID) (uploadPath, mimeType string, err error) {
	db := psql.DBContext(ctx)

	var attachment models.Attachment

	if err := db.Model(&models.Attachment{}).First(&attachment, "id = ?", attachmentID).Error; err != nil {
		return "", "", errs.NotFoundErr().WithError(err)
	}

	return a.fileService.GetPath(attachment.UploadPath), attachment.MimeType, nil
}
