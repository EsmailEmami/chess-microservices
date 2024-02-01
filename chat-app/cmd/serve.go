package cmd

import (
	"github.com/esmailemami/chess/chat/internal/app/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		server.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
