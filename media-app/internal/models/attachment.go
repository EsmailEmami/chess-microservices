package models

import (
	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
)

type Attachment struct {
	models.Model

	MimeType     string    `gorm:"column:mime_type"                 json:"mimeType"`
	Extension    string    `gorm:"column:extension"                 json:"extension"`
	OriginalName string    `gorm:"column:original_name"             json:"originalName"`
	FileName     string    `gorm:"column:file_name"                 json:"fineName"`
	FileType     string    `gorm:"column:file_type"                 json:"fileType"`
	UploadPath   string    `gorm:"column:upload_path"               json:"uploadPath"`
	ItemID       uuid.UUID `gorm:"column:item_id"                   json:"itemId"`
}

func (Attachment) TableName() string {
	return "media.attachment"
}

const (
	ATTACHMENT_USER_PROFILE        = "user-profile"
	ATTACHMENT_PUBLIC_ROOM_PROFILE = "public-room-profile"
	ATTACHMENT_ROOM_FILE_MESSAGE   = "room-file-message"
)
