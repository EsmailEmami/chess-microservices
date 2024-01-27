package models

import (
	"github.com/esmailemami/chess/chat/internal/models"
	baseconsts "github.com/esmailemami/chess/shared/consts"
	sharedModels "github.com/esmailemami/chess/shared/models"
	"github.com/esmailemami/chess/shared/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CreateGlobalRoomInputModel struct {
	Name  string      `json:"name"`
	Users []uuid.UUID `json:"users"`
}

func (model CreateGlobalRoomInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.Name,
			validation.Required.Error(baseconsts.Required),
		),
		validation.Field(
			&model.Users,
			validation.Each(validation.By(validations.ExistsInDB(&sharedModels.User{}, "id", baseconsts.RecordNotFound))),
		),
	)
}

func (c *CreateGlobalRoomInputModel) ToDBModel() *models.Room {
	r := &models.Room{
		Name:      c.Name,
		IsPrivate: false,
	}
	r.ID = uuid.New()

	return r
}

type CreatePrivateRoomInputModel struct {
	UserID uuid.UUID `json:"userId"`
}

func (model CreatePrivateRoomInputModel) Validate() error {
	return validation.ValidateStruct(
		&model,
		validation.Field(
			&model.UserID,
			validation.Required.Error(baseconsts.Required),
			validation.By(validations.ExistsInDB(&sharedModels.User{}, "id", baseconsts.RecordNotFound)),
		),
	)
}

func (c *CreatePrivateRoomInputModel) ToDBModel() *models.Room {
	r := &models.Room{
		Name:      "private room",
		IsPrivate: true,
	}
	r.ID = uuid.New()

	return r
}

// Single room output

type RoomOutPutModell struct {
	ID        uuid.UUID `gorm:"id;type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"name" json:"name"`
	IsPrivate bool      `gorm:"is_private" json:"isPrivate"`
	Users     []User    `gorm:"foreignKey:room_id" json:"users"`
	Messages  []Message `gorm:"foreignKey:room_id" json:"messages"`
}

type User struct {
	ID     uuid.UUID   `gorm:"id;type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID   `gorm:"user_id;type:uuid" json:"userId"`
	RoomID uuid.UUID   `gorm:"room_id;type:uuid" json:"roomId"`
	User   UserProfile `gorm:"foreignKey:user_id;references:id;" json:"user"`
}

type UserProfile struct {
	ID        uuid.UUID `gorm:"id;type:uuid;primaryKey" json:"id"`
	FirstName *string   `gorm:"first_name" json:"firstName"`
	LastName  *string   `gorm:"last_name" json:"lastName"`
	Username  string    `gorm:"username" json:"username"`
}

type Message struct {
	ID        uuid.UUID  `gorm:"id;type:uuid;primaryKey" json:"id"`
	Content   string     `gorm:"content" json:"content"`
	ReplyToID *uuid.UUID `gorm:"reply_to_id;type:uuid" json:"replyToId"`
	ReplyTo   *Message   `gorm:"foreignKey:ReplyToID;references:ID" json:"replyTo"`
	RoomID    uuid.UUID  `gorm:"room_id;type:uuid" json:"roomId"`
}

type RoomOutPutModel struct {
	ID        uuid.UUID `gorm:"id;type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"name" json:"name"`
	IsPrivate bool      `gorm:"is_private" json:"isPrivate"`
	Users     []struct {
		ID     uuid.UUID `gorm:"id;type:uuid;primaryKey" json:"id"`
		UserID uuid.UUID `gorm:"user_id;type:uuid" json:"userId"`
		RoomID uuid.UUID `gorm:"room_id;type:uuid" json:"roomId"`
		User   struct {
			ID        uuid.UUID `gorm:"id;type:uuid;primaryKey" json:"id"`
			FirstName *string   `gorm:"first_name" json:"firstName"`
			LastName  *string   `gorm:"last_name" json:"lastName"`
			Username  string    `gorm:"username" json:"username"`
		} `gorm:"foreignKey:user_id;references:id;" json:"user"`
	} `gorm:"foreignKey:room_id" json:"users"`
	Messages []struct {
		ID        uuid.UUID  `gorm:"id;type:uuid;primaryKey" json:"id"`
		Content   string     `gorm:"content" json:"content"`
		ReplyToID *uuid.UUID `gorm:"reply_to_id;type:uuid" json:"replyToId"`
		ReplyTo   *struct {
			ID      uuid.UUID `gorm:"id;type:uuid;primaryKey" json:"id"`
			Content string    `gorm:"content" json:"content"`
		} `gorm:"foreignKey:reply_to_id;references:id" json:"replyTo"`
		RoomID uuid.UUID `gorm:"room_id;type:uuid" json:"roomId"`
	} `gorm:"foreignKey:room_id" json:"messages"`
}
