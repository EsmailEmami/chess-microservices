package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/chat/internal/app/service"

	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

var (
	PublicRoomProfileChangedCh = make(chan *RoomProfileMessage, 256)
)

const (
	mediaExchange = "media_user_ex"
)

func initializeMediaRabbitMQ() {
	go consumeMediaRoomProfile()
}

func consumeMediaRoomProfile() {
	msgsCh, err := amqp.ConsumeMessagesFromExchange("media_chat_queue", mediaExchange, "media_chat.room_profile")

	if err != nil {
		logging.ErrorE("failed to consume rabbit MQ", err)
		return
	}

	roomService := service.NewRoomService(redis.GetConnection())

	for msg := range msgsCh {
		var resp RoomProfileMessage

		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			logging.ErrorE("failed to unmarshal room_profile consumer", err)
		}

		if err := roomService.UpdateAvatar(context.Background(), resp.RoomID, resp.ProfilePath); err != nil {
			logging.Error("Failed to set room avatar", "roomId", resp.RoomID)
			continue
		}

		PublicRoomProfileChangedCh <- &resp
	}
}

type RoomProfileMessage struct {
	RoomID      uuid.UUID `json:"roomId"`
	ProfilePath string    `json:"profilePath"`
}
