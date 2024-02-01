package chatroom

import (
	"github.com/esmailemami/chess/chat/internal/app/models"
	"github.com/google/uuid"
)

// global chat room message
type RoomMessage struct {
	RoomID uuid.UUID `json:"roomId"`
	Data   any       `json:"data"`
}

type RoomOutPutModel struct {
	Room     *models.RoomOutPutModel   `json:"room"`
	Messages []models.MessageOutPutDto `json:"messages"`
}
