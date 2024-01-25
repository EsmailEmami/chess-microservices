package psql

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DefaultConfig = &gorm.Config{
	FullSaveAssociations: false,
}

var gormDBConn *gorm.DB

func DBContext(ctx context.Context) *gorm.DB {
	if gormDBConn == nil {
		panic("database is not initialized")
	}

	return gormDBConn.WithContext(ctx)
}

func Initialize(user, password, host, dbName, port, sslmode string, config *gorm.Config) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s  TimeZone=Asia/Tehran",
		host, port, user, dbName, password, sslmode,
	)

	dbConn, err := gorm.Open(postgres.Open(connStr), config)
	if err != nil {
		return err
	}
	dbConn = dbConn.Omit(clause.Associations)

	loadCallbacks(dbConn)

	sqlDB, _ := dbConn.DB()
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	gormDBConn = dbConn

	return nil
}
