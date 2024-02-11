package service

import (
	"context"
	"mime/multipart"

	"github.com/esmailemami/chess/media/internal/models"
	"github.com/esmailemami/chess/media/internal/util"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
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

// UploadFile must delete the last one and replace the new image
func (a *AttachmentService) UploadFile(ctx context.Context, f *multipart.FileHeader, itemID uuid.UUID, directory string, fileType string) (*models.Attachment, error) {
	db := psql.DBContext(ctx)
	tx := db.Begin()

	var attachment models.Attachment

	if err := db.Model(&models.Attachment{}).Order("created_at DESC").Find(&attachment, "item_id = ?", itemID).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	newAttachment, err := a.uploadFileTransaction(tx, f, itemID, directory, fileType)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if attachment.ID != uuid.Nil {

		if err := a.fileService.DeleteFile(attachment.UploadPath); err != nil {
			_ = a.fileService.DeleteFile(newAttachment.UploadPath)
			tx.Rollback()
			return nil, errs.InternalServerErr().WithError(err)
		}

		tx.Delete(&attachment)
	}

	tx.Commit()

	return newAttachment, nil
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
