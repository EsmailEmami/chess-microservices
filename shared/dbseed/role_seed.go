package dbseed

import (
	"github.com/esmailemami/chess/shared/models"
	"gorm.io/gorm"
)

func seedRole(dbConn *gorm.DB) error {
	items := []models.Role{
		{
			Model: models.Model{
				ID: models.ROLE_ROOT,
			},
			Name:     "root",
			Code:     "1",
			IsSystem: true,
		},
		{
			Model: models.Model{
				ID: models.ROLE_ADMIN,
			},
			Name:     "admin",
			Code:     "2",
			IsSystem: true,
		},
		{
			Model: models.Model{
				ID: models.ROLE_USER,
			},
			Name:     "user",
			Code:     "3",
			IsSystem: true,
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
