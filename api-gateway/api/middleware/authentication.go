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
		if r.Header.Get("UserId") == "" {
			errs.ErrorHandler(w, errs.UnAuthorizedErr())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := util.GetTokenString(r)

		if err == nil {
			authService := service.GetAuthGrpcClient()

			resp, err := authService.Authenticate(r.Context(), &grpc.AuthenticateRequest{
				Token: token,
			})

			if err == nil {
				r.Header.Add("UserId", resp.UserId)
				r.Header.Add("JwtId", resp.UserId)
			}
		}

		next.ServeHTTP(w, r)
	})
}
