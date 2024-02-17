package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/google/uuid"
)

const (
	userExchange = "media_user_events"
)

func initializeUserRabbitMQ() {
	// initialize the exchange between 'user' and 'media' apps
	if err := amqp.DeclareExchange(userExchange, rabbitmq.Topic, true, false); err != nil {
		logging.FatalE("failed to declare 'media_user_events' exchange", err)
	}
}

func PublishUserProfile(ctx context.Context, attachmentID, userID uuid.UUID, profilePath string, deleteAttachmentID *uuid.UUID) error {
	body := &userProfileMessage{
		UserID:      userID,
		ProfilePath: profilePath,
	}

	correlationId := ""

	if deleteAttachmentID != nil {
		correlationId = deleteAttachmentID.String()
	}

	return amqp.PublishMessage(ctx, userExchange, "media_user.profile.upload", &rabbitmq.Message{
		Body:          body,
		DeliveryMode:  rabbitmq.DeliveryModePersistent,
		CorrelationId: correlationId,
		ReplyTo:       "media_callbacks.user.profile.delete",
		AppId:         "media-app",
	})
}

type userProfileMessage struct {
	UserID      uuid.UUID `json:"userId"`
	ProfilePath string    `json:"profilePath"`
}
