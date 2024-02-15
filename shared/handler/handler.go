package handler

import (
	"mime/multipart"

	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/models"
	"github.com/gin-gonic/gin"
)

const (
	TenMB = 10 << 20
)

type Handler struct {
}

func (*Handler) GetUser(c *gin.Context) *models.User {
	return c.Value("user").(*models.User)
}

func (*Handler) GetFiles(c *gin.Context, maximumSize int64) ([]*multipart.FileHeader, error) {
	err := c.Request.ParseMultipartForm(maximumSize)
	if err != nil {
		return nil, errs.BadRequestErr().WithError(err).Msg("Invalid file size")
	}

	var files []*multipart.FileHeader

	for _, fileHeaders := range c.Request.MultipartForm.File {
		files = append(files, fileHeaders...)
	}

	return files, err
}
