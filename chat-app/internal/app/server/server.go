package server

import (
	"log"

	"github.com/esmailemami/chess/chat/api/routes"
	"github.com/esmailemami/chess/chat/docs"
	chatroom "github.com/esmailemami/chess/chat/internal/chat-room"
	"github.com/esmailemami/chess/chat/internal/websocket"
	"github.com/esmailemami/chess/chat/pkg/rabbitmq"
	"github.com/esmailemami/chess/shared/consul"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Wellcome to auth service"))
	})

	setupSwagger(r)

	routes.Initialize(r)
	websocket.InitializeRoutes(r)

	rabbitmq.Initialize()

	chatroom.Run()

	go websocket.Run()

	go consul.Register()

	log.Fatal(r.Run(":" + viper.GetString("app.port")))
}

func setupSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "Chat API doc"
	docs.SwaggerInfo.Description = "Chat API."
	docs.SwaggerInfo.Version = "1.0"
	url := viper.GetString("app.url")

	docs.SwaggerInfo.Host = url
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
