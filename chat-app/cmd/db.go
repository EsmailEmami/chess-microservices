package cmd

import (
	"context"

	"github.com/esmailemami/chess/chat/internal/seed"
	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/dbseed"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use: "db",
}

func init() {
	rootCmd.AddCommand(dbCmd)

	dbCmd.AddCommand(&cobra.Command{
		Use: "seed",
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbseed.Run(psql.DBContext(context.Background()), seed.SeedRoom)
		},
	})
}
