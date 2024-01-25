package dbseed

import (
	"time"

	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedUser(dbConn *gorm.DB) error {
	items := []models.User{
		{
			Model: models.Model{
				ID: uuid.MustParse("85fb2028-787f-4102-8464-503a9e362f57"),
			},
			FirstName: func() *string {
				value := "esmail"
				return &value
			}(),
			LastName: func() *string {
				value := "emami"
				return &value
			}(),
			Username: "esmailemami",
			Password: "$2a$10$2oV2MylgwZftP47vL/ndteC6tzmcY85qRNo/5FTCeS403eL8zo9Yq",
			Mobile: func() *string {
				value := "09903669556"
				return &value
			}(),
			Enabled: true,
			RoleID:  models.ROLE_ROOT,
		},
		{
			Model: models.Model{
				ID: uuid.MustParse("a59e4d0c-ddd5-46bc-befb-eb36dcc13eea"),
			},
			FirstName: func() *string {
				value := "esmail"
				return &value
			}(),
			LastName: func() *string {
				value := "emami"
				return &value
			}(),
			Username: "esmailemami2",
			Password: "$2a$10$2oV2MylgwZftP47vL/ndteC6tzmcY85qRNo/5FTCeS403eL8zo9Yq",
			Mobile: func() *string {
				value := "09903669556"
				return &value
			}(),
			Enabled: true,
			RoleID:  models.ROLE_ADMIN,
		},
	}

	for _, item := range items {

		var old models.User
		err := dbConn.Where("id", item.ID).First(&old).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		if err != nil && err == gorm.ErrRecordNotFound {
			err = dbConn.Save(&item).Error
			if err != nil {
				return err
			}
		} else {
			err = dbConn.Model(&models.User{}).Where("id", item.ID).UpdateColumns(map[string]any{
				"first_name": item.FirstName,
				"last_name":  item.LastName,
				"username":   item.Username,
				"mobile":     item.Mobile,
				"enabled":    item.Enabled,
				"role_id":    item.RoleID.String(),
				"updated_at": time.Now(),
			}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
