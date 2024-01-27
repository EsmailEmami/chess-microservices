package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageOutPutDto struct {
	ID        uuid.UUID `gorm:"id" json:"id"`
	Content   string    `gorm:"content" json:"content"`
	CreatedAt time.Time `gorm:"created_at" json:"createdAt"`

	UserID    *uuid.UUID `gorm:"created_by_id" json:"userId"`
	FirstName *string    `gorm:"first_name" json:"firstName"`
	LastName  *string    `gorm:"last_name" json:"lastName"`

	ReplyToID        *uuid.UUID `gorm:"reply_to_id" json:"replyToId"`
	ReplyToContent   *string    `gorm:"reply_content" json:"replyToContent"`
	ReplyToFirstName *string    `gorm:"reply_first_name" json:"replyToFirstName"`
	ReplyToLastName  *string    `gorm:"reply_last_name" json:"replyToLastName"`
}
