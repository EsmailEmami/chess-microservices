package middleware

import (
	"net/http"

	"github.com/esmailemami/chess/api-gateway/api/util"
	"github.com/esmailemami/chess/api-gateway/internal/grpc"
	"github.com/esmailemami/chess/api-gateway/internal/service"
	"github.com/esmailemami/chess/shared/errs"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := util.GetTokenString(r)
		if err != nil {
			errs.ErrorHandler(w, errs.UnAuthorizedErr())
			return
		}

		authService := service.GetAuthGrpcClient()

		resp, err := authService.Authenticate(r.Context(), &grpc.AuthenticateRequest{
			Token: token,
		})

		if err != nil {
			errs.ErrorHandler(w, errs.UnAuthorizedErr())
			return
		}

		r.Header.Add("UserId", resp.UserId)
		next.ServeHTTP(w, r)
	})
}
