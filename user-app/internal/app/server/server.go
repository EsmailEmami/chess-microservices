package server

import (
	"log"

	"github.com/esmailemami/chess/shared/consul"
	"github.com/esmailemami/chess/user/api/routes"
	"github.com/esmailemami/chess/user/docs"
	consumerRMQ "github.com/esmailemami/chess/user/internal/rabbitmq"
	producerRMQ "github.com/esmailemami/chess/user/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Wellcome to user service"))
	})

	routes.Initialize(r)

	setupSwagger(r)

	producerRMQ.InitializeProducerConnection()
	consumerRMQ.InitializeConsumerConnection()

	go consul.Register()

	port := viper.GetString("app.port")
	log.Fatal(r.Run(":" + port))
}

func setupSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "User API doc"
	docs.SwaggerInfo.Description = "User API."
	docs.SwaggerInfo.Version = "1.0"
	url := viper.GetString("app.url")

	docs.SwaggerInfo.Host = url
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
