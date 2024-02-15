package rabbitmq

import (
	"context"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

const (
	userExchange = "chat_user_ex"
)

var (
	UserProfileChangedCh = make(chan *UserProfileChangedMessage, 256)
)

func initializeUserRabbitMQ() {
	amqp.DeclareExchange(userExchange, rabbitmq.Direct)

	userLastConnectionQueue, err := amqp.DeclareQueue("chat_user_last_connection_queue", true, false, false)

	if err != nil {
		logging.FatalE("failed to declare 'chat_user_queue' queue", err)
	}

	if err := amqp.BindQueueToExchange(userLastConnectionQueue.Name, userExchange, "chat_user_last_connection"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	go consumeUserProfile()
}

func PublishUserLastConnection(ctx context.Context, userID uuid.UUID, date time.Time) error {
	return amqp.PublishMessage(ctx, userExchange, "chat_user_last_connection", &struct {
		UserID uuid.UUID `json:"userId"`
		Date   time.Time `json:"date"`
	}{
		UserID: userID,
		Date:   date,
	})
}

func consumeUserProfile() {
	msgsCh, err := amqp.ConsumeMessagesFromExchange("chat_user_profile_queue", userExchange, "chat_user_profile")

	if err != nil {
		logging.ErrorE("failed to consume rabbit MQ", err)
		return
	}

	logging.Debug("Chat Profile Receivning....")

	for msg := range msgsCh {
		var data UserProfileChangedMessage
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			logging.ErrorE("failed to unmarshal user profile consumer", err)
		}

		logging.Info("Chat Profile Received")

		UserProfileChangedCh <- &data
	}
}

type UserProfileChangedMessage struct {
	UserID      uuid.UUID `json:"userID"`
	ProfilePath string    `json:"profilePath"`
}
