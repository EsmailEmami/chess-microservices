package service

import (
	"github.com/esmailemami/chess/chat/internal/models"
	"github.com/esmailemami/chess/shared/service"
)

type RoomService struct {
	service.BaseService[models.Room]
}

func NewRoomService() *RoomService {
	return &RoomService{}
}
