package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/google/uuid"
)

const (
	chatExchange = "media_chat_ex"
)

func initializeChatRabbitMQ() {
	amqp.DeclareExchange(chatExchange, rabbitmq.Direct)

	queue, err := amqp.DeclareQueue("media_chat_profile_queue", true, false, false)
	if err != nil {
		logging.FatalE("failed to declare 'media_chat_profile_queue' queue", err)
	}

	if err := amqp.BindQueueToExchange(queue.Name, chatExchange, "media_chat_profile"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}
}

func PublishRoomProfile(ctx context.Context, roomID uuid.UUID, profilePath string) error {
	return amqp.PublishMessage(ctx, chatExchange, "media_chat_profile", &roomProfileMessage{
		RoomID:      roomID,
		ProfilePath: profilePath,
	})
}

type roomProfileMessage struct {
	RoomID      uuid.UUID `json:"roomId"`
	ProfilePath string    `json:"profilePath"`
}
