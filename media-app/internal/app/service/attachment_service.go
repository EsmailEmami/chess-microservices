package service

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/esmailemami/chess/media/internal/models"
	"github.com/esmailemami/chess/media/internal/util"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AttachmentService struct {
	directory string
}

func NewAttachmentService() *AttachmentService {
	return &AttachmentService{
		directory: viper.GetString("app.upload_files_directory"),
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
				_ = os.Remove(attachment.UploadPath)
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
	uploadPath, fileName, err := a.uploadFile(f, directory)

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
		_ = os.Remove(uploadPath)

		return nil, errs.InternalServerErr().WithError(err)
	}

	return attachment, err
}

func (a *AttachmentService) uploadFile(f *multipart.FileHeader, directory string) (uploadPath, fileName string, err error) {
	bts, err := a.readFile(f)
	if err != nil {
		return "", "", errs.InternalServerErr().WithError(err)
	}

	absDirectory, err := a.getPath(directory)
	if err != nil {
		return "", "", errs.InternalServerErr().WithError(err)
	}

	uploadPath, fileName = util.GenerateFilename(absDirectory, f.Filename)

	err = a.writeFile(uploadPath, bts)
	if err != nil {
		return "", "", err
	}

	return uploadPath, fileName, nil
}

func (a *AttachmentService) readFile(f *multipart.FileHeader) ([]byte, error) {
	file, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func (a *AttachmentService) getPath(directory string) (string, error) {
	path := path.Join(a.directory, directory)

	// make sure directory exists
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (a *AttachmentService) writeFile(filePath string, b []byte) error {
	fileToWrite, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	defer fileToWrite.Close()
	_, err = fileToWrite.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = fileToWrite.Write(b)
	if err != nil {
		return err
	}
	return fileToWrite.Sync()
}
