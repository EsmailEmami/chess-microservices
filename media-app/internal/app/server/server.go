package server

import (
	"log"
	"net/http"

	"github.com/esmailemami/chess/media/api/routes"
	"github.com/esmailemami/chess/media/docs"
	"github.com/esmailemami/chess/media/pkg/rabbitmq"
	"github.com/esmailemami/chess/shared/consul"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Wellcome to media service"))
	})

	setupSwagger(r)

	routes.Initialize(r)

	rabbitmq.Initialize()

	go consul.Register()

	filesDir := viper.GetString("app.upload_files_directory")

	r.StaticFS("/uploads", http.Dir(filesDir))

	log.Fatal(r.Run(":" + viper.GetString("app.port")))
}

func setupSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "Media API doc"
	docs.SwaggerInfo.Description = "Media API."
	docs.SwaggerInfo.Version = "1.0"
	url := viper.GetString("app.url")

	docs.SwaggerInfo.Host = url
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
