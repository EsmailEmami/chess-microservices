package rabbitmq

import (
	"context"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

const (
	chatExchange = "chat_user_ex"
)

func initializeChatRabbitMQ() {
	go consumeLastConnection()
}

type lastConnectionMessage struct {
	UserID uuid.UUID `json:"userId"`
	Date   time.Time `json:"date"`
}

func consumeLastConnection() {
	msgsCh, err := amqp.ConsumeMessagesFromExchange("chat_user_last_connection_queue", chatExchange, "chat_user_last_connection")

	if err != nil {
		logging.ErrorE("failed to consume rabbit MQ", err)
		return
	}

	userService := service.NewUserService()

	for msg := range msgsCh {
		var data lastConnectionMessage
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			logging.ErrorE("failed to unmarshal last connection consumer", err)
		}

		if err := userService.UpdateLastConnection(context.Background(), data.UserID, data.Date); err != nil {
			logging.Error("failed to update user connection", "userId", data.UserID)
		}
	}
}

func PublishUserProfile(ctx context.Context, userID uuid.UUID, profilePath string) error {
	logging.Info("User Profile Published")

	return amqp.PublishMessage(ctx, chatExchange, "chat_user_profile", &userProfileMessage{
		UserID:      userID,
		ProfilePath: profilePath,
	})
}
