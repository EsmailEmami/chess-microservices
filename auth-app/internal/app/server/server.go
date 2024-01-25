package server

import (
	"fmt"
	"log"
	"net"

	appGrpc "github.com/esmailemami/chess/auth/internal/app/grpc"
	"github.com/esmailemami/chess/auth/internal/app/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RunGrpcServer() {
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
