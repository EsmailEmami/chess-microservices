package dbutil

import (
	"context"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/logging"
)

func Exists(model interface{}, condition interface{}, args ...interface{}) bool {
	db := psql.DBContext(context.Background())

	var count int64
	if err := db.Model(model).Where(condition, args...).Count(&count).Error; err != nil {
		logging.ErrorE("db exists failed to execute", err)
		return false
	}

	return count > 0
}
