package service

import (
	appGrpc "github.com/esmailemami/chess/api-gateway/internal/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var authServiceClient appGrpc.AuthServiceClient

func GetAuthGrpcConnection() (*grpc.ClientConn, error) {
	return grpc.Dial(":"+viper.GetString("grpc.auth_port"), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func InitializeServices(conn *grpc.ClientConn) {
	authServiceClient = appGrpc.NewAuthServiceClient(conn)
}

func GetAuthGrpcClient() appGrpc.AuthServiceClient {
	if authServiceClient != nil {
		return authServiceClient
	}

	panic("service not implemented")
}
