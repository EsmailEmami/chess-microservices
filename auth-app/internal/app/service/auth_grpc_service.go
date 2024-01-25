package service

import (
	"context"

	"github.com/esmailemami/chess/auth/internal/app/grpc"
	"github.com/google/uuid"
)

type AuthGrpcService struct {
	grpc.UnimplementedAuthServiceServer

	authService *AuthService
}

func NewAuthGrpcService(authService *AuthService) grpc.AuthServiceServer {
	return &AuthGrpcService{
		authService: authService,
	}
}

func (a *AuthGrpcService) Authenticate(ctx context.Context, req *grpc.AuthenticateRequest) (*grpc.AuthenticateResponse, error) {
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

	return &grpc.AuthenticateResponse{
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
	}, nil
}
