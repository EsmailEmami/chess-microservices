package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/esmailemami/chess/user/internal/util"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

const (
	mediaExchange = "media_user_ex"
)

func initializeMediaRabbitMQ() {
	go consumeMediaUserProfile()
}

func consumeMediaUserProfile() {
	msgsCh, err := amqp.ConsumeMessagesFromExchange("media_user_profile_queue", mediaExchange, "media_user_profile")

	if err != nil {
		logging.ErrorE("failed to consume rabbit MQ", err)
		return
	}

	userService := service.NewUserService()

	for msg := range msgsCh {
		var resp userProfileMessage

		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			logging.ErrorE("failed to unmarshal profile consumer", err)
		}

		if err := userService.UpdateProfile(context.Background(), resp.UserID, resp.ProfilePath); err != nil {
			logging.Error("Failed to set user profile", "userId", resp.UserID)
		}

		if err := PublishUserProfile(context.Background(), resp.UserID, util.FilePathPrefix(resp.ProfilePath)); err != nil {
			logging.Error("Failed to publish user profile", "userId", resp.UserID)
		}
	}
}

type userProfileMessage struct {
	UserID      uuid.UUID `json:"userId"`
	ProfilePath string    `json:"profilePath"`
}
