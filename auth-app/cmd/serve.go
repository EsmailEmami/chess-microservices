package cmd

import (
	"log"

	"github.com/esmailemami/chess/auth/api/routes"
	"github.com/esmailemami/chess/auth/internal/app/server"
	"github.com/esmailemami/chess/shared/consul"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/", func(ctx *gin.Context) {
			ctx.Writer.Write([]byte("Wellcome to auth service"))
		})

		// initialize the routes
		routes.Initialize(r)

		// run grpc
		go server.RunGrpcServer()

		// register consul
		go consul.Register()

		log.Fatal(r.Run(":" + viper.GetString("app.port")))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
