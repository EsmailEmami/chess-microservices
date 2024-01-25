package migrations

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

//go:embed *yaml
var migrationFS embed.FS

func init() {
	log.SetFlags(log.Lshortfile)
}

type migrationYamlFile struct {
	Up         string `yaml:"up"`
	Down       string `yaml:"down"`
	Connection string `yaml:"connection"`
}

type migration struct {
	ID       uint
	Name     string
	Batch    uint
	Filename string
	App      string
}

func (m migration) Down(dbConn *gorm.DB) error {
	name := m.Name
	fmt.Println("Rollback ", name)
	filename := m.Filename

	bts, err := migrationFS.ReadFile(filename)
	if err != nil {
		return err
	}

	var mf migrationYamlFile
	err = yaml.Unmarshal(bts, &mf)
	if err != nil {
		return err
	}
	migrateSql := mf.Down

	var tx *gorm.DB

	// handle connections
	switch mf.Connection {
	case "log-db":
		return nil
	default:
		tx = dbConn.Begin()
	}

	err = tx.Exec(migrateSql).Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	err = dbConn.Delete(&m).Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// MakeMigration create new migration file
func MakeMigration(dirPath, create string) error {
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%v/%v_%v.yaml", dirPath, timestamp, strcase.ToSnake(create))

	content := `---
up: |
  -- UP SQL

down: |
  -- DOWN SQL
`
	f, e := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
	if e != nil {
		fmt.Println(e)
		return e
	}
	defer f.Close()
	_, err := f.WriteString(content)
	return err
}

// checkMigrationTable will check existence of migrations table and create it if it doesn't exist.
func checkMigrationTable(dbConn *gorm.DB) {
	dbConn.Exec(`CREATE SCHEMA IF NOT EXISTS public;`)
	dbConn.Exec(`CREATE TABLE IF NOT EXISTS migrations (
    id serial NOT NULL,
    name varchar(255) NOT NULL,
    batch integer,
    filename varchar(512) NOT NULL,
	app varchar(512) NOT NULL,
    CONSTRAINT migrations_pkey PRIMARY KEY (id)
    )`)
}

// Migrate call Migrate function of files and save batches in database
func Migrate(dbConn *gorm.DB, app, migrationsPath string) error {
	checkMigrationTable(dbConn)

	// global migrations
	files, e := migrationFS.ReadDir(".")
	if e != nil {
		return errors.New("Error on loading local migrations directory")
	}

	// app migrations

	appMigs, e := os.ReadDir(migrationsPath)

	if e != nil {
		return errors.New("Error on loading application migrations directory")
	}

	files = append(files, appMigs...)

	var batch struct {
		LastBatch uint `sql:"last_batch"`
	}
	dbConn.Raw(`select max(batch) as last_batch from migrations;`).Scan(&batch)
	batch.LastBatch++

	upToDate := true

	for _, v := range files {
		if v.IsDir() {
			continue
		}
		var mg migration

		filename := v.Name()
		migrationName := strings.TrimSuffix(filename, ".yml")
		migrationName = strings.TrimSuffix(migrationName, ".yaml")
		underscoreIndex := strings.Index(filename, "_")
		migrationName = migrationName[underscoreIndex+1:]

		_ = dbConn.Where("name LIKE ?", migrationName).First(&mg).Error

		if mg.ID != 0 {
			//This file already migrated
			continue
		}

		var mf migrationYamlFile

		var fReader io.Reader

		file, err := migrationFS.Open(v.Name())
		// maybe it is not global file
		if err == nil {
			fReader = file
		} else {
			bts, err := os.ReadFile(path.Join(migrationsPath, filename))

			if err != nil {
				return fmt.Errorf("could not read file %s, %+v", filename, err)
			}

			fReader = bytes.NewReader(bts)
		}

		err = yaml.NewDecoder(fReader).Decode(&mf)
		if err != nil {
			return err
		}

		var tx *gorm.DB

		migrateSql := mf.Up
		switch mf.Connection {
		case "log-db":
			continue
		default:
			tx = dbConn.Begin()
		}

		err = tx.Exec(migrateSql).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}

		// save migrated file to DB
		mg.Name = migrationName
		mg.Batch = batch.LastBatch
		mg.Filename = filename
		batch.LastBatch = batch.LastBatch + 1
		mg.App = app
		err = dbConn.Create(&mg).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}

		tx.Commit()

		fmt.Println("Migrated ", migrationName)
		upToDate = false
	}

	if upToDate {
		fmt.Println("Already up to date!")
	}

	return dbConn.Exec("update migrations set batch=coalesce((select max(batch) from migrations) , 0)+1 where batch is null;").Error
}

// Rollback will rollback database using batch number
func Rollback(dbConn *gorm.DB, app string) error {
	checkMigrationTable(dbConn)

	var rollbacks []migration
	dbConn.Where("batch = (select max(batch) from migrations where app = ?)", app).Find(&rollbacks)
	for _, v := range rollbacks {
		err := v.Down(dbConn)
		if err != nil {
			return err
		}
	}

	return nil
}

func RollbackAll(dbConn *gorm.DB, app string) error {
	checkMigrationTable(dbConn)

	var rollbacks []migration
	dbConn.Order("id desc").Where("app = ?", app).Find(&rollbacks)
	for _, v := range rollbacks {
		err := v.Down(dbConn)
		if err != nil {
			return err
		}
	}

	return nil
}
