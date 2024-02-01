package server

import (
	"log"

	"github.com/esmailemami/chess/shared/consul"
	"github.com/esmailemami/chess/user/api/routes"
	"github.com/esmailemami/chess/user/internal/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Wellcome to game service"))
	})

	routes.Initialize(r)

	rabbitmq.Initialize()

	go consul.Register()

	port := viper.GetString("app.port")
	log.Fatal(r.Run(":" + port))
}
