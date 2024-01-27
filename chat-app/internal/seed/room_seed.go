package seed

import (
	"github.com/esmailemami/chess/chat/internal/models"
	sharedModels "github.com/esmailemami/chess/shared/models"
	"gorm.io/gorm"
)

func SeedRoom(dbConn *gorm.DB) error {
	items := []models.Room{
		{
			Model: sharedModels.Model{
				ID: models.GlobalRoomID,
			},
			Name:      "global chat room",
			IsPrivate: false,
		},
	}

	for _, item := range items {
		err := dbConn.Where("id", item.ID).FirstOrCreate(&item).Error
		if err != nil {
			return err
		}
	}

	return nil
}
