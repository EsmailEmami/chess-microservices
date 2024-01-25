package handler

import (
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (*Handler) GetUser(c *gin.Context) (*models.User, error) {
	user := c.Value("user").(*models.User)

	if user == nil {
		return nil, errs.NotFoundErr()
	}
	return user, nil
}
