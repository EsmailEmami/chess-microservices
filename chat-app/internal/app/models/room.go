package models

import (
	"github.com/esmailemami/chess/chat/internal/models"
	baseconsts "github.com/esmailemami/chess/shared/consts"
	sharedModels "github.com/esmailemami/chess/shared/models"
	"github.com/esmailemami/chess/shared/validations"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CreatePublicRoomInputModel struct {
	Name  string      `json:"name"`
	Users []uuid.UUID `json:"users"`
}

func (model CreatePublicRoomInputModel) Validate() error {
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

func (c *CreatePublicRoomInputModel) ToDBModel() *models.Room {
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
		IsPrivate: true,
	}
	r.ID = uuid.New()

	return r
}

type RoomOutPutModel struct {
	ID        uuid.UUID             `json:"id"`
	Name      string                `json:"name"`
	IsPrivate bool                  `json:"isPrivate"`
	Users     []RoomUserOutPutModel `json:"users"`
}

type RoomUserOutPutModel struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Username  string    `json:"username"`
}

type RoomQueryParams struct {
	SearchTerm string `json:"searchTerm"`
	SortColumn string `json:"sortColumn"`
	SortOrder  string `json:"sortOrder"`
	Page       int    `json:"page" default:"1"`
	Limit      int    `json:"limit" default:"25"`
}

type RoomsOutPutModel struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
