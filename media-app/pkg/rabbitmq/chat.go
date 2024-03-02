package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/google/uuid"
)

const (
	chatExchange = "media_chat_events"
)

func initializeChatRabbitMQ() {
	// initialize the exchange between 'user' and 'media' apps
	if err := amqp.DeclareExchange(chatExchange, rabbitmq.Topic, true, false); err != nil {
		logging.FatalE("failed to declare 'media_chat_events' exchange", err)
	}
}

func PublishRoomAvatar(ctx context.Context, attachmentID, roomID uuid.UUID, filePath string, deleteAttachmentID *uuid.UUID) error {
	body := &roomProfileMessage{
		RoomID:      roomID,
		ProfilePath: filePath,
	}

	correlationId := ""

	if deleteAttachmentID != nil {
		correlationId = deleteAttachmentID.String()
	}

	return amqp.PublishMessage(ctx, chatExchange, "media_chat.room.profile.upload", &rabbitmq.Message{
		Body:          body,
		DeliveryMode:  rabbitmq.DeliveryModePersistent,
		CorrelationId: correlationId,
		ReplyTo:       "media_callbacks.chat.room.profile.delete",
		AppId:         "media-app",
	})
}

type roomProfileMessage struct {
	RoomID      uuid.UUID `json:"roomId"`
	ProfilePath string    `json:"profilePath"`
}

func PublishRoomFileMessage(ctx context.Context, roomID, messageID, userID uuid.UUID, fileID string, fileType string) error {
	body := &roomFileMessage{
		RoomID:    roomID,
		UserID:    userID,
		MessageID: messageID,
		Type:      fileType,
		File:      fileID,
	}

	return amqp.PublishMessage(ctx, chatExchange, "media_chat.room.message.file.upload", &rabbitmq.Message{
		Body:         body,
		DeliveryMode: rabbitmq.DeliveryModePersistent,
		AppId:        "media-app",
	})
}

type roomFileMessage struct {
	RoomID    uuid.UUID `json:"roomId"`
	UserID    uuid.UUID `json:"userId"`
	MessageID uuid.UUID `json:"messageId"`
	Type      string    `json:"type"`
	File      string    `json:"file"`
}
