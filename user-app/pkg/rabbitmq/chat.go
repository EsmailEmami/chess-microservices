package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/google/uuid"
)

const (
	chatExchange = "chat_user_events"
)

func initializeChatRabbitMQ() {
	// initialize the exchange between 'user' and 'media' apps
	if err := amqp.DeclareExchange(chatExchange, rabbitmq.Topic, true, false); err != nil {
		logging.FatalE("failed to declare 'media_user_events' exchange", err)
	}
}

func PublishChatRoomUserProfileChangedMessage(ctx context.Context, userID uuid.UUID, profilePath string) error {
	body := &userProfileMessage{
		UserID:      userID,
		ProfilePath: profilePath,
	}

	return amqp.PublishMessage(ctx, chatExchange, "chat_user.user.profile.changed", &rabbitmq.Message{
		Body: body,
	})
}

type userProfileMessage struct {
	UserID      uuid.UUID `json:"userId"`
	ProfilePath string    `json:"profilePath"`
}
