package util

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

func GetMimeType(f *multipart.FileHeader) string {
	return f.Header.Get("Content-Type")
}

func GetFileExetension(filename string) string {
	p := strings.Split(filename, ".")
	if len(p) <= 1 {
		return ""
	}

	return strings.ToLower(p[len(p)-1])
}

func GenerateFilename(directory, orginalFilename string) (uploadPath, fileName string) {
	ext := GetFileExetension(orginalFilename)

	fileName = fmt.Sprintf("%s.%s", uuid.New(), ext)
	uploadPath = fmt.Sprintf("%s/%s", directory, fileName)

	if _, err := os.Stat(uploadPath); err == nil {
		return uploadPath, fileName
	}

	return GenerateFilename(directory, orginalFilename)
}
