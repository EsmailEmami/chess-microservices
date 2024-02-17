package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/esmailemami/chess/user/internal/util"
	"github.com/esmailemami/chess/user/pkg/rabbitmq"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

const (
	mediaExchange = "media_user_events"
)

func initializeMediaRabbitMQ() {
	go consumeMediaUserProfileUpload()
}

func consumeMediaUserProfileUpload() {
	userProfileQueue, err := amqp.DeclareQueue("media_user_profile_upload", true, false, false)
	if err != nil {
		logging.FatalE("failed to declare 'media_user_profile_upload' queue", err)
	}

	if err := amqp.BindQueueToExchange(userProfileQueue.Name, mediaExchange, "media_user.profile.upload"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	messageBus, err := amqp.ConsumeMessages(userProfileQueue.Name, false)
	if err != nil {
		logging.FatalE("failed to consume 'media_user_profile_upload' queue", err)
	}

	userService := service.NewUserService()

	for msg := range messageBus {
		logging.Debug("user profile upload message received")

		var resp userProfileMessage

		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			logging.ErrorE("failed to unmarshal profile consumer", err)
		}

		if err := userService.UpdateProfile(context.Background(), resp.UserID, resp.ProfilePath); err != nil {
			logging.ErrorE("Failed to set user profile", err, "userId", resp.UserID)
		}

		// if err := PublishUserProfile(context.Background(), resp.UserID, util.FilePathPrefix(resp.ProfilePath)); err != nil {
		// 	logging.Error("Failed to publish user profile", "userId", resp.UserID)
		// }

		if err := msg.Ack(false); err != nil {
			logging.ErrorE("failed to acknowlendge 'media_user_profile_upload' queue", err)
			continue
		}

		// send the new user profile to chat rooms
		rabbitmq.PublishChatRoomUserProfileChangedMessage(context.Background(), resp.UserID, util.FilePathPrefix(resp.ProfilePath))

		// send the callback to the media app for deleting the last profile
		rabbitmq.PublishUserProfieUploadedMessage(context.Background(), msg.ReplyTo, msg.CorrelationId)
	}
}

type userProfileMessage struct {
	UserID      uuid.UUID `json:"userId"`
	ProfilePath string    `json:"profilePath"`
}
