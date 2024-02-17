package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
)

const (
	mediaExchange = "media_callbacks"
)

func PublishRoomAvatarUploadedMessage(ctx context.Context, replyTo, correlationId string) error {
	return amqp.PublishMessage(ctx, mediaExchange, replyTo, &rabbitmq.Message{
		CorrelationId: correlationId,
		Body: &successAcknowledgeMessage{
			Message: "Hell yeah!",
		},
	})
}

type successAcknowledgeMessage struct {
	Message string `json:"message"`
}
