package models

import "github.com/esmailemami/chess/shared/models"

type Attachment struct {
	models.Model

	MimeType     string `gorm:"column:mime_type"                 json:"mimeType"`
	Extension    string `gorm:"column:extension"                 json:"extension"`
	OriginalName string `gorm:"column:original_name"             json:"originalName"`
	FileName     string `gorm:"column:file_name"                 json:"fineName"`
	FileType     string `gorm:"column:file_type"                 json:"fileType"`
}
