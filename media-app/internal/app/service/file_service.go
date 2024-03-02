package service

import (
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/esmailemami/chess/media/internal/util"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/spf13/viper"
)

type FileService struct {
	directory string
}

func NewFileService() *FileService {
	return &FileService{
		directory: viper.GetString("app.upload_files_directory"),
	}
}

func (f *FileService) UploadFile(file *multipart.FileHeader, directory string) (uploadPath, fileName string, err error) {
	bts, err := f.readFile(file)
	if err != nil {
		return "", "", errs.InternalServerErr().WithError(err)
	}

	absDirectory, err := f.getDirectory(directory)
	if err != nil {
		return "", "", errs.InternalServerErr().WithError(err)
	}

	fileName = util.GenerateFilename(absDirectory, file.Filename)

	err = f.writeFile(path.Join(absDirectory, fileName), bts)
	if err != nil {
		return "", "", err
	}

	uploadPath = path.Join(directory, fileName)
	return uploadPath, fileName, nil
}

func (f *FileService) readFile(file *multipart.FileHeader) ([]byte, error) {
	fileOpen, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileOpen.Close()

	return io.ReadAll(fileOpen)
}

func (f *FileService) getDirectory(directory string) (string, error) {
	path := f.GetPath(directory)

	// make sure directory exists
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (f *FileService) fileExists(filePath string) bool {
	path := f.GetPath(filePath)
	return util.FileExists(path)
}

func (f *FileService) writeFile(filePath string, b []byte) error {
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

func (f *FileService) GetPath(filePath string) string {
	return path.Join(f.directory, filePath)
}

func (f *FileService) DeleteFile(path string) error {
	if !f.fileExists(path) {
		return nil
	}

	absPath := f.GetPath(path)

	return os.Remove(absPath)
}
