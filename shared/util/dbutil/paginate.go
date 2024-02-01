package dbutil

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/esmailemami/chess/shared/errs"
	"gorm.io/gorm"
)

func dbLikeSearch(searchTerm string) (string, bool) {
	if strings.TrimSpace(searchTerm) != "" {
		return "%" + strings.TrimSpace(searchTerm) + "%", true
	}
	return "", false
}

func Filter(db *gorm.DB, searchTerm string, fields ...string) *gorm.DB {
	if searchTerm, ok := dbLikeSearch(searchTerm); ok && len(fields) > 0 {
		columns := make([]string, len(fields))
		values := make([]interface{}, len(fields))

		for i, column := range fields {
			fmt.Println("column:", column)

			columns[i] = column + " ILIKE ?"
			values[i] = searchTerm
		}

		db = db.Where(strings.Join(columns, " OR "), values...)
	}

	return db
}

func Paginate(db *gorm.DB, page, limit int, data any) (totalRecords int64, err error) {
	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		if err := db.WithContext(context.Background()).Count(&totalRecords).Error; err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		qry := db.WithContext(context.Background())

		qry = paginate(qry, page, limit)
		if err := qry.Find(data).Error; err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		return 0, errs.InternalServerErr().WithError(err)
	default:
	}

	return
}

func paginate(db *gorm.DB, page, limit int) *gorm.DB {
	return db.Offset(limit * (page - 1)).Limit(limit)
}
