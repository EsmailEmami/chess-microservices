package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/google/uuid"
)

const (
	userExchange = "media_user_ex"
)

func initializeUserRabbitMQ() {
	amqp.DeclareExchange(userExchange, rabbitmq.Topic)

	queue, err := amqp.DeclareQueue("media_user_queue", true, false, false)

	if err != nil {
		logging.FatalE("failed to declare 'media_user_queue' queue", err)
	}

	if err := amqp.BindQueueToExchange(queue.Name, userExchange, "media_user.*"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}
}

func PublishUserProfile(ctx context.Context, userID uuid.UUID, profilePath string) error {
	return amqp.PublishMessage(ctx, userExchange, "media_user.profile", &userProfileMessage{
		UserID:      userID,
		ProfilePath: profilePath,
	})
}

type userProfileMessage struct {
	UserID      uuid.UUID `json:"userId"`
	ProfilePath string    `json:"profilePath"`
}
