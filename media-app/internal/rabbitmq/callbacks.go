package rabbitmq

import (
	"context"
	"strings"
	"time"

	"github.com/esmailemami/chess/media/internal/app/service"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func initializeCallbacks() {
	if err := amqp.DeclareExchange("media_callbacks", rabbitmq.Topic, true, false); err != nil {
		logging.FatalE("failed to declare 'media_user_events' exchange", err)
	}

	go consumeUserProfileAck()
	go consumeChatRoomAvatarAck()
}

func consumeUserProfileAck() {
	userProfileQueue, err := amqp.DeclareQueue("media_user_profile_ack", true, false, false)
	if err != nil {
		logging.FatalE("failed to declare 'media_user_profile_ack' queue", err)
	}

	if err := amqp.BindQueueToExchange(userProfileQueue.Name, "media_callbacks", "media_callbacks.user.profile.#"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	messagesBus, err := amqp.ConsumeMessages(userProfileQueue.Name, true)

	if err != nil {
		logging.FatalE("failed to consume 'media_user_profile_ack' queue", err)
	}

	attachmentService := service.NewAttachmentService()

	g := errgroup.Group{}
	g.SetLimit(10)

	for msg := range messagesBus {
		attachmentID, err := uuid.Parse(msg.CorrelationId)

		if err != nil {
			logging.ErrorE("failed to parse attachment id", err)
			continue
		}

		g.Go(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			err := attachmentService.Delete(ctx, attachmentID)
			if err != nil {
				logging.ErrorE("failed to delete attachment", err)
			} else {
				logging.Debug("file deleted successfuly")
			}
			return err
		})
	}
}

func consumeChatRoomAvatarAck() {
	userProfileQueue, err := amqp.DeclareQueue("media_chat_room_avatar_ack", true, false, false)
	if err != nil {
		logging.FatalE("failed to declare 'media_chat_room_avatar_ack' queue", err)
	}

	if err := amqp.BindQueueToExchange(userProfileQueue.Name, "media_callbacks", "media_callbacks.chat.room.profile.#"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	messagesBus, err := amqp.ConsumeMessages(userProfileQueue.Name, true)

	if err != nil {
		logging.FatalE("failed to consume 'media_chat_room_avatar_ack' queue", err)
	}

	attachmentService := service.NewAttachmentService()

	g := errgroup.Group{}
	g.SetLimit(10)

	for msg := range messagesBus {
		if strings.TrimSpace(msg.CorrelationId) == "" {
			continue
		}

		attachmentID, err := uuid.Parse(msg.CorrelationId)

		if err != nil {
			logging.ErrorE("failed to parse attachment id", err)
			continue
		}

		g.Go(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			err := attachmentService.Delete(ctx, attachmentID)
			if err != nil {
				logging.ErrorE("failed to delete attachment", err)
			} else {
				logging.Debug("file deleted successfuly")
			}
			return err
		})
	}
}
