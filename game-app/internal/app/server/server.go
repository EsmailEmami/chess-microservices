package server

import (
	"log"

	"github.com/esmailemami/chess/game/api/routes"
	"github.com/esmailemami/chess/game/internal/app/chess"
	"github.com/esmailemami/chess/game/pkg/websocket"
	"github.com/esmailemami/chess/shared/consul"
	"github.com/esmailemami/chess/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Wellcome to game service"))
	})

	// initialize the routes
	routes.Initialize(r)

	// run the websockets
	websocket.Run()

	// run chess game
	go chess.Run()

	// register consul
	go consul.Register()

	ws := r.Group("/ws")
	ws.Use(middleware.Authorization())
	ws.GET("/chess", websocket.ChessWss.HandleWS)

	port := viper.GetString("app.port")
	log.Fatal(r.Run(":" + port))
}
