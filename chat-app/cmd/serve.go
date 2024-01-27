package cmd

import (
	"log"

	"github.com/esmailemami/chess/chat/api/routes"
	globalroom "github.com/esmailemami/chess/chat/internal/global-room"
	"github.com/esmailemami/chess/chat/pkg/websocket"
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

		// global game room
		go globalroom.Run()

		// run websocket
		go websocket.Run()

		websocket.InitializeRoutes(r)

		// register consul
		go consul.Register()

		log.Fatal(r.Run(":" + viper.GetString("app.port")))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
