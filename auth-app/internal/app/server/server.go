package server

import (
	"fmt"
	"log"
	"net"

	"github.com/esmailemami/chess/auth/api/routes"
	"github.com/esmailemami/chess/auth/docs"
	appGrpc "github.com/esmailemami/chess/auth/internal/app/grpc"
	"github.com/esmailemami/chess/auth/internal/app/service"
	"github.com/esmailemami/chess/shared/consul"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Wellcome to auth service"))
	})

	// initialize the routes
	routes.Initialize(r)

	setupSwagger(r)

	// run grpc
	go runGrpcServer()

	// register consul
	go consul.Register()

	log.Fatal(r.Run(":" + viper.GetString("app.port")))
}

func runGrpcServer() {
	port := viper.GetString("app.grpc_port")

	net, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer net.Close()

	var (
		userService = service.NewUserService()
		authService = service.NewAuthService(userService)
	)
	server := grpc.NewServer()
	appGrpc.RegisterAuthServiceServer(server, service.NewAuthGrpcService(authService))

	fmt.Println("grpc started on port", port)
	if err := server.Serve(net); err != nil {
		log.Fatal(err)
	}
}

func setupSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "Auth API doc"
	docs.SwaggerInfo.Description = "Auth API."
	docs.SwaggerInfo.Version = "1.0"
	url := viper.GetString("app.url")

	docs.SwaggerInfo.Host = url
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
