package grpc

import (
	"context"

	"github.com/esmailemami/chess/auth/internal/app/service"
	"github.com/esmailemami/chess/auth/proto"

	"github.com/google/uuid"
)

type AuthGrpcService struct {
	proto.UnimplementedAuthServiceServer

	authService *service.AuthService
}

func NewAuthGrpcService(authService *service.AuthService) proto.AuthServiceServer {
	return &AuthGrpcService{
		authService: authService,
	}
}

func (a *AuthGrpcService) Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	token, err := a.authService.ParseTokenString(req.Token, true)
	if err != nil {
		return nil, err
	}

	tokenID, err := uuid.Parse(token.JwtID())
	if err != nil {
		return nil, err
	}

	user, err := a.authService.GetUser(ctx, tokenID)

	if err != nil {
		return nil, err
	}

	return &proto.AuthenticateResponse{
		UserId:   user.ID.String(),
		Username: user.Username,
		FirstName: func() string {
			if user.FirstName != nil {
				return *user.FirstName
			}
			return ""
		}(),
		LastName: func() string {
			if user.LastName != nil {
				return *user.LastName
			}
			return ""
		}(),
		JwtId: token.JwtID(),
	}, nil
}
