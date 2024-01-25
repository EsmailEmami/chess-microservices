package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/migrations"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var migrationsPath string

const app = "user"

// migrationCmd represents the migration command
var migrationCmd = &cobra.Command{
	Use: "migration",
}

func init() {
	checkDir()

	rootCmd.AddCommand(migrationCmd)

	migrationCmd.AddCommand(&cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrations.Migrate(migrationDB(), app, migrationsPath)
		},
	})

	migrationCmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrations.Rollback(migrationDB(), app)
		},
	})

	migrationCmd.AddCommand(&cobra.Command{
		Use:   "make",
		Short: "Making a new migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("migration name is required")
			}
			name := args[0]
			return migrations.MakeMigration(migrationsPath, name)
		},
	})

	migrationCmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Rollback all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrations.RollbackAll(migrationDB(), app)
		},
	})
}

func checkDir() {
	err := os.MkdirAll("../migrations", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	migrationsPath, err = filepath.Abs("../migrations")

	if err != nil {
		log.Fatal(err)
	}
}

func migrationDB() *gorm.DB {
	return psql.DBContext(context.Background()).Session(&gorm.Session{SkipHooks: false, NewDB: false, Logger: logger.Discard})
}
