package validations

import (
	"context"
	"errors"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/util"
	"github.com/google/uuid"
)

func ExistsInDB(model interface{}, column string, errorMsg string) func(value interface{}) error {
	return func(value interface{}) error {
		if util.IsNil(value) {
			return nil
		}

		db := psql.DBContext(context.Background())

		var count int64
		db.Model(model).
			Where(column+"=?", value).
			Count(&count)

		if count > 0 {
			return nil
		}

		return errors.New(errorMsg)
	}
}

func ExistsInDBWithCond(model interface{}, column string, errorMsg string, condition interface{}, args ...interface{}) func(value interface{}) error {
	return func(value interface{}) error {
		if util.IsNil(value) {
			return nil
		}

		db := psql.DBContext(context.Background())

		var count int64
		db.Model(model).
			Where(column+"=?", value).
			Where(condition, args...).
			Count(&count)

		if count > 0 {
			return nil
		}

		return errors.New(errorMsg)
	}
}

func NotExistsInDB(model interface{}, column string, errorMsg string, id ...uuid.UUID) func(value interface{}) error {
	return func(value interface{}) error {
		if util.IsNil(value) {
			return nil
		}

		db := psql.DBContext(context.Background())

		var count int64
		checkDB := db.Model(model).
			Where(column+"=?", value)

		if len(id) > 0 {
			checkDB = checkDB.Where("id != ", id[0])
		}

		checkDB.Count(&count)

		if count > 0 {
			errs.InternalServerErr()

			return errors.New(errorMsg)
		}

		return nil
	}
}

func NotExistsInDBWithCond(model interface{}, column string, errorMsg string, condition interface{}, args ...interface{}) func(value interface{}) error {
	return func(value interface{}) error {
		if util.IsNil(value) {
			return nil
		}

		db := psql.DBContext(context.Background())

		var count int64
		db.Model(model).
			Where(column+"=?", value).
			Where(condition, args...).
			Count(&count)

		if count > 0 {
			return errors.New(errorMsg)
		}

		return nil
	}
}
