package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageOutPutDto struct {
	ID        uuid.UUID `gorm:"column:id" json:"id"`
	Content   string    `gorm:"column:content" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`

	UserID    uuid.UUID `gorm:"column:user_id" json:"userId"`
	FirstName *string   `gorm:"column:first_name" json:"firstName"`
	LastName  *string   `gorm:"last_name" json:"lastName"`
	Type      string    `gorm:"column:type" json:"type"`
	IsPin     bool      `gorm:"column:is_pin" json:"isPin"`

	ReplyToID        *uuid.UUID `gorm:"column:reply_to_id" json:"replyToId"`
	ReplyToContent   *string    `gorm:"column:reply_content" json:"replyToContent"`
	ReplyToFirstName *string    `gorm:"column:reply_first_name" json:"replyToFirstName"`
	ReplyToLastName  *string    `gorm:"column:reply_last_name" json:"replyToLastName"`

	IsEdited bool `gorm:"column:is_edited" json:"isEdited"`
	IsSeen   bool `gorm:"column:is_seen"   json:"isSeen"`
}
