package rabbitmq

import (
	"context"

	"github.com/esmailemami/chess/chat/internal/app/service"
	"github.com/esmailemami/chess/chat/pkg/rabbitmq"
	"github.com/esmailemami/chess/shared/database/redis"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

var (
	PublicRoomProfileChangedCh = make(chan *RoomAvatarMessage, 256)
	RoomFileMessageCh          = make(chan *RoomFileMessage, 256)
)

const (
	mediaExchange = "media_chat_events"
)

func initializeMediaRabbitMQ() {
	go consumeMediaRoomAvatarUpload()
	go consumeMediaFileMessage()
}

func consumeMediaRoomAvatarUpload() {
	userProfileQueue, err := amqp.DeclareQueue("media_chat_room_avatar_upload", true, false, false)
	if err != nil {
		logging.FatalE("failed to declare 'media_chat_room_avatar_upload' queue", err)
	}

	if err := amqp.BindQueueToExchange(userProfileQueue.Name, mediaExchange, "media_chat.room.profile.upload"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	messageBus, err := amqp.ConsumeMessages(userProfileQueue.Name, false)
	if err != nil {
		logging.FatalE("failed to consume 'media_chat_room_avatar_upload' queue", err)
	}

	var (
		cache          = redis.GetConnection()
		messageService = service.NewMessageService(cache)
		roomService    = service.NewRoomService(cache, messageService)
	)

	for msg := range messageBus {
		logging.Debug("room avatar upload message received")

		var resp RoomAvatarMessage

		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			logging.ErrorE("failed to unmarshal profile consumer", err)
		}

		if err := roomService.UpdateAvatar(context.Background(), resp.RoomID, resp.ProfilePath); err != nil {
			logging.ErrorE("Failed to update room avatar", err, "id", resp.RoomID)
		}
		if err := msg.Ack(false); err != nil {
			logging.ErrorE("failed to acknowlendge 'media_chat_room_avatar_upload' queue", err)
			continue
		}

		PublicRoomProfileChangedCh <- &resp

		// send the callback to the media app for deleting the last profile
		rabbitmq.PublishRoomAvatarUploadedMessage(context.Background(), msg.ReplyTo, msg.CorrelationId)
	}
}

type RoomAvatarMessage struct {
	RoomID      uuid.UUID `json:"roomId"`
	ProfilePath string    `json:"profilePath"`
}

func consumeMediaFileMessage() {
	queue, err := amqp.DeclareQueue("media_chat_room_file_message_upload", true, false, false)
	if err != nil {
		logging.FatalE("failed to declare 'media_chat_room_file_message_upload' queue", err)
	}

	if err := amqp.BindQueueToExchange(queue.Name, mediaExchange, "media_chat.room.message.file.upload"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	messageBus, err := amqp.ConsumeMessages(queue.Name, false)
	if err != nil {
		logging.FatalE("failed to consume 'media_chat_room_file_message_upload' queue", err)
	}

	for msg := range messageBus {
		logging.Debug("room avatar upload message received")

		var resp RoomFileMessage

		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			logging.ErrorE("failed to unmarshal file message consumer", err)
		}

		if err := msg.Ack(false); err != nil {
			logging.ErrorE("failed to acknowlendge 'media_chat_room_avatar_upload' queue", err)
			continue
		}

		RoomFileMessageCh <- &resp
	}
}

type RoomFileMessage struct {
	RoomID    uuid.UUID `json:"roomId"`
	UserID    uuid.UUID `json:"userId"`
	MessageID uuid.UUID `json:"messageId"`
	Type      string    `json:"type"`
	File      string    `json:"file"`
}
