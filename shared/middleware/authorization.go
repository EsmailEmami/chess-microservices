package middleware

import (
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.GetHeader("UserId")
		if userIDStr == "" {
			ctx.Abort()
			errs.ErrorHandler(ctx.Writer, errs.UnAuthorizedErr())
			return
		}

		userID, err := uuid.Parse(userIDStr)

		if err != nil {
			ctx.Abort()
			errs.ErrorHandler(ctx.Writer, errs.UnAuthorizedErr())
			return
		}

		userService := service.NewUserService()

		user, err := userService.Get(ctx, userID)

		if err != nil {
			ctx.Abort()
			errs.ErrorHandler(ctx.Writer, errs.UnAuthorizedErr())
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
